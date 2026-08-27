package nagad

// InitializeRequest represents the data sent to start a payment
type InitializeRequest struct {
	AccountNo string `json:"accountNumber"`
	DateTime  string `json:"datetime"`
}

// InitializeResponse represents the response from Nagad initialization
type InitializeResponse struct {
	SensitiveData string `json:"sensitiveData"`
	Signature     string `json:"signature"`
}

// DecryptedSensitiveData represents the decrypted content from Nagad's initialization response
type DecryptedSensitiveData struct {
	PaymentReferenceID string `json:"paymentReferenceId"`
	PlainKey           string `json:"plainKey"`
	IV                 string `json:"iv"`
}

// CompleteRequest represents the data sent to finalize a payment
type CompleteRequest struct {
	SensitiveData       string `json:"sensitiveData"`
	Signature           string `json:"signature"`
	MerchantCallbackURL string `json:"merchantCallbackURL"`
}

// CompleteResponse represents the response from Nagad completion
type CompleteResponse struct {
	CallBackURL string `json:"callBackUrl"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
}

// PaymentDetails is the inner data that gets AES encrypted
type PaymentDetails struct {
	MerchantID string `json:"merchantId"`
	OrderID    string `json:"orderId"`
	Amount     string `json:"amount"`
	Currency   string `json:"currencyCode"`
	Challenge  string `json:"challenge"`
}

// VerifyResponse represents the transaction status response
type VerifyResponse struct {
	MerchantID         string `json:"merchantId"`
	OrderID            string `json:"orderId"`
	PaymentReferenceID string `json:"paymentReferenceId"`
	Amount             string `json:"amount"`
	Status             string `json:"status"`
	StatusCode         string `json:"statusCode"`
	DateTime           string `json:"datetime"`
	IssuerPaymentRef   string `json:"issuerPaymentReference"`
}
