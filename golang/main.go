package main

import (
	"encoding/json"
	"fmt"
	"github.com/Subskribe/samples/service"
	"github.com/Subskribe/samples/usage"
	"os"
	"time"
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

	log.Infof("Headers: %s", respTenants.Header)
	log.Infof("Body: %s", string(respTenants.Body))
	log.Infof("")

	respUsers, err := svc.Get("/users")

	if err != nil {
		log.Errorf("Error getting users %s", err)
		return
	}
	log.Infof("Headers: %s", respUsers.Header)
	log.Infof("Body: %s", string(respUsers.Body))

	usageData := &usage.RawUsageData{
		Data: []usage.RawUsage{
			{UsageId: "usg-001", AliasId: "als-001", UsageTime: time.Now().Unix(), UsageQuantity: 100},
			{UsageId: "usg-001", SubscriptionId: "SUB-001", ChargeId: "CHRG-001", UsageTime: time.Now().Unix(), UsageQuantity: 100},
		},
	}

	payload, err := json.Marshal(usageData)
	if err != nil {
		log.Errorf("Cannot marshall usage data: %s", err)
		return
	}
	log.Infof("about to submit usage payload: %s", string(payload))
	respUsage, err := svc.Post("/v2/usage", payload, service.JsonContentType)

	if err != nil {
		log.Errorf("Error calling usage submision: %s", err)
		return
	}
	log.Infof("Headers: %s", respUsage.Header)
	log.Infof("Body: %s", string(respUsage.Body))
}
