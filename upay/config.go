package upay

import "net/http"

type Environment string

const (
	Sandbox    Environment = "sandbox"
	Production Environment = "production"
)

type Config struct {
	MerchantID           string
	MerchantKey          string
	MerchantCode         string
	MerchantName         string
	MerchantMobile       string
	MerchantCity         string
	MerchantCategoryCode string
	Environment          Environment
	HTTPClient           *http.Client
}

func (c *Config) GetBaseURL() string {
	if c.Environment == Production {
		return "https://api.upaybd.com"
	}
	return "https://pg-sandbox.upaybd.com" // Assuming sandbox URL
}
