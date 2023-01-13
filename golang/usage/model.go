package usage

type RawUsage struct {
	//UsageId idempotent identifier for usage resubmitting with same id does nothing
	UsageId        string `json:"id"`
	//AliasId is a friendly identifier linking to the right subscription and charge
	AliasId        string `json:"aliasId"`
	//SubscriptionId the subscription for which this usage is for
	SubscriptionId string `json:"subscriptionId"`
	//ChargeId the charge for which this usage is for
	ChargeId       string `json:"chargeId"`
	//UsageTime should be in unix epoch time format
	UsageTime     int64 `json:"usageTime"`
	//UsageQuantity the usage quantity
	UsageQuantity int64 `json:"usageQuantity"`
}

type RawUsageData struct {
	//Data slice of raw usage
	Data []RawUsage `json:"data"`
}
