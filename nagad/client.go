package nagad

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	config *Config
}

func NewClient(config *Config) *Client {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &Client{config: config}
}

func (c *Client) Initialize(orderID string, clientIP string) (*DecryptedSensitiveData, error) {
	url := fmt.Sprintf("%scheck-out/initialize/%s/%s", c.config.GetBaseURL(), c.config.MerchantID, orderID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req, clientIP)

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("initialization failed with status %d: %s", resp.StatusCode, string(body))
	}

	var initResp InitializeResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
		return nil, err
	}

	// Decrypt sensitive data to get paymentReferenceId, plainKey, and iv
	decryptedData, err := DecryptRSA(c.config.MerchantPrivateKey, initResp.SensitiveData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt sensitive data: %w", err)
	}

	var sensitive DecryptedSensitiveData
	if err := json.Unmarshal(decryptedData, &sensitive); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sensitive data: %w", err)
	}

	return &sensitive, nil
}

func (c *Client) Complete(paymentRefID string, details PaymentDetails, sensitiveData DecryptedSensitiveData, clientIP string, callbackURL string) (*CompleteResponse, error) {
	// 1. AES Encrypt PaymentDetails using plainKey and iv
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	encryptedDetails, err := AESEncrypt([]byte(sensitiveData.PlainKey), []byte(sensitiveData.IV), detailsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt payment details: %w", err)
	}

	// 2. Sign the encrypted details
	signature, err := SignRSA(c.config.MerchantPrivateKey, detailsJSON) // Nagad docs say sign the PLAIN JSON
	if err != nil {
		return nil, fmt.Errorf("failed to sign payment details: %w", err)
	}

	completeReq := CompleteRequest{
		SensitiveData:       encryptedDetails,
		Signature:           signature,
		MerchantCallbackURL: callbackURL,
	}

	url := fmt.Sprintf("%scheck-out/complete/%s", c.config.GetBaseURL(), paymentRefID)
	body, _ := json.Marshal(completeReq)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	c.setHeaders(req, clientIP)

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var compResp CompleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&compResp); err != nil {
		return nil, err
	}

	return &compResp, nil
}

func (c *Client) Verify(paymentRefID string) (*VerifyResponse, error) {
	url := fmt.Sprintf("%sverify/payment/%s", c.config.GetBaseURL(), paymentRefID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var verifyResp VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return nil, err
	}

	return &verifyResp, nil
}

func (c *Client) setHeaders(req *http.Request, clientIP string) {
	req.Header.Set(HeaderAPIKey, c.config.APIKey)
	req.Header.Set(HeaderIPV4, clientIP)
	req.Header.Set(HeaderClientType, ClientTypePCWeb)
}
