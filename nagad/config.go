package nagad

import "net/http"

// Environment represents the Nagad API environment (Sandbox or Production)
type Environment string

const (
	Sandbox    Environment = "sandbox"
	Production Environment = "production"
)

// Config holds the merchant credentials and environment settings
type Config struct {
	MerchantID         string
	MerchantPrivateKey string
	NagadPublicKey     string
	APIKey             string
	Environment        Environment
	HTTPClient         *http.Client
}

// GetBaseURL returns the Nagad API base URL for the configured environment
func (c *Config) GetBaseURL() string {
	if c.Environment == Production {
		return "https://api.mynagad.com/api/dfs/"
	}
	return "http://sandbox.mynagad.com:10080/remote-payment-gateway-1.0/api/dfs/"
}
