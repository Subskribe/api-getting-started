package service

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

//DefaultTimeout is set to 5 seconds for any HTTP call
const DefaultTimeout = time.Second * 5

//ApiKeyHeader is the authentication authorization header to call the API in Subskribe
const ApiKeyHeader = "X-API-Key"
const ContentTypeHeader = "Content-Type"
//LocationHeader is used to extract the identifiers out API response in case of POST requests
const LocationHeader = "Location"
const JsonContentType = "application/json"

var emptyParamMap = map[string]string{}

type HttpService interface {
	Get(path string) (*Response, error)

	GetQuery(path string, params map[string]string) (*Response, error)

	MultiPartPost(path string, params map[string]string, paramName, filePath string) (*Response, error)

	Put(path string, body []byte) (*Response, error)

	Post(path string, body []byte, contentType string) (*Response, error)

	Delete(path string) (*Response, error)
}

type service struct {
	baseUrl string
	apiKey  string
	timeout time.Duration
}

type Response struct {
	Body   []byte
	Header http.Header
}

func (s *service) makeClient() *http.Client {
	return &http.Client{Timeout: s.timeout}
}

func (s *service) newSvcRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", s.baseUrl, path), body)
	if err != nil {
		return nil, err
	}
	header := make([]string, 1)
	header[0] = s.apiKey
	req.Header[ApiKeyHeader] = header
	return req, nil
}

func (s *service) Get(path string) (*Response, error) {
	return s.GetQuery(path, emptyParamMap)
}

func (s *service) GetQuery(path string, params map[string]string) (*Response, error) {
	client := s.makeClient()
	req, err := s.newSvcRequest(http.MethodGet, path, nil)

	if err != nil {
		return nil, err
	}

	// deal with query params
	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	return s.svcDo(req, client)
}

func (s *service) Put(path string, body []byte) (*Response, error) {
	client := s.makeClient()
	req, err := s.newSvcRequest(http.MethodPut, path, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}
	return s.svcDo(req, client)
}

func (s *service) Post(path string, body []byte, contentType string) (*Response, error) {
	client := s.makeClient()
	req, err := s.newSvcRequest(http.MethodPost, path, bytes.NewReader(body))

	if contentType != "" {
		req.Header.Set(ContentTypeHeader, contentType)
	}

	if err != nil {
		return nil, err
	}
	return s.svcDo(req, client)
}

func (s *service) Delete(path string) (*Response, error) {
	client := s.makeClient()
	// TODO: assuming deletes are nil body requests provide body arg if needed
	req, err := s.newSvcRequest(http.MethodDelete, path, nil)

	if err != nil {
		return nil, err
	}
	return s.svcDo(req, client)
}

func (s *service) svcDo(req *http.Request, client *http.Client) (*Response, error) {
	log.Infof("request URL: %s", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return processResponse(resp)
}

func (s *service) MultiPartPost(path string, params map[string]string, paramName, filePath string) (*Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fpart, err := writer.CreateFormFile(paramName, filepath.Base(filePath))

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fpart, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		err = writer.WriteField(key, val)
		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := s.newSvcRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(ContentTypeHeader, writer.FormDataContentType())

	client := s.makeClient()
	return s.svcDo(req, client)
}

func processResponse(resp *http.Response) (*Response, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return &Response{Body: body, Header: resp.Header}, nil
	}

	return nil, fmt.Errorf("http request failed status:%d response:%s", resp.StatusCode, string(body))
}

//NewService creates a new service instance
func NewService(baseUrl, apiKey string, timeout time.Duration) HttpService {
	return &service{baseUrl: baseUrl, apiKey: apiKey, timeout: timeout}
}

func LocationKey(header http.Header) (string, error) {
	loc, ok := header[LocationHeader]

	if !ok {
		return "", fmt.Errorf("could not find %s in headers", LocationHeader)
	}

	key := path.Base(loc[0])
	if key == "" {
		return "", fmt.Errorf("cloud not find location in: %s", loc[0])
	}
	return key, nil
}