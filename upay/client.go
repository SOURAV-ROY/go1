package upay

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
	token  string
}

func NewClient(config *Config) *Client {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &Client{config: config}
}

func (c *Client) Auth() error {
	url := fmt.Sprintf("%s/api/v1/merchant/auth", c.config.GetBaseURL())
	authReq := AuthRequest{
		MerchantID:  c.config.MerchantID,
		MerchantKey: c.config.MerchantKey,
	}
	body, _ := json.Marshal(authReq)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	c.token = authResp.Data.Token
	return nil
}

func (c *Client) CreatePayment(invoiceID string, amount float64, redirectURL string) (*CreatePaymentResponse, error) {
	if c.token == "" {
		if err := c.Auth(); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/api/v1/payment/create", c.config.GetBaseURL())
	createReq := CreatePaymentRequest{
		MerchantID:              c.config.MerchantID,
		MerchantKey:             c.config.MerchantKey,
		MerchantCode:            c.config.MerchantCode,
		MerchantName:            c.config.MerchantName,
		MerchantMobile:          c.config.MerchantMobile,
		MerchantCity:            c.config.MerchantCity,
		MerchantCategoryCode:    c.config.MerchantCategoryCode,
		TransactionCurrencyCode: "BDT",
		Amount:                  amount,
		InvoiceID:               invoiceID,
		TxnID:                   fmt.Sprintf("TXN-%s", invoiceID),
		Date:                    time.Now().Format("2006-01-02"),
		RedirectURL:             redirectURL,
	}

	body, _ := json.Marshal(createReq)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createResp CreatePaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, err
	}

	return &createResp, nil
}

func (c *Client) QueryPayment(invoiceID string) (*QueryPaymentResponse, error) {
	if c.token == "" {
		if err := c.Auth(); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/api/v1/payment/query", c.config.GetBaseURL())
	queryReq := QueryPaymentRequest{
		MerchantID: c.config.MerchantID,
		InvoiceID:  invoiceID,
	}

	body, _ := json.Marshal(queryReq)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var queryResp QueryPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&queryResp); err != nil {
		return nil, err
	}

	return &queryResp, nil
}
