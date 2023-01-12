package main

import (
	"fmt"
	"github.com/Subskribe/samples/service"
	"os"
)
import log "github.com/sirupsen/logrus"

const apiKeyEnvVariable = "SUBSKRIBE_API_KEY"
const apiBaseUrlEnvVariable = "SUBSKRIBE_BASE_URL"
// NOTE: please do not include a trailing slash in the base url
const defaultApiBaseUrl = "https://billy-sandbox.subskribe.net"

func main() {
	fmt.Println("Starting API session with Subskribe Service")

	// get API key and based URL
	apiKey, ok := os.LookupEnv(apiKeyEnvVariable)
	if !ok {
		log.Errorf("could not find api key env variable: %s", apiKeyEnvVariable)
		return
	}

	baseUrl, ok := os.LookupEnv(apiBaseUrlEnvVariable)

	if !ok {
		baseUrl = defaultApiBaseUrl
	}

	svc := service.NewService(baseUrl, apiKey, service.DefaultTimeout)

	// tenant info urls
	respTenants, err := svc.Get("/tenants")
	if err != nil {
		log.Errorf("Error getting tenant information %s", err)
		return
	}

	log.Infof("Headers: %s \n", respTenants.Header)
	log.Infof("Headers: %s \n", string(respTenants.Body))
	log.Infof("")

	respUsers, err := svc.Get("/users")

	if err != nil {
		log.Errorf("Error getting users %s", err)
		return
	}
	log.Infof("Headers: %s \n", respUsers.Header)
	log.Infof("Headers: %s \n", string(respUsers.Body))
}
