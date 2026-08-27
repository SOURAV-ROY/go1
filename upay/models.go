package upay

type AuthRequest struct {
	MerchantID  string `json:"merchant_id"`
	MerchantKey string `json:"merchant_key"`
}

type AuthResponse struct {
	Code    string `json:"code"`
	Lang    string `json:"lang"`
	Message string `json:"message"`
	Data    struct {
		MerchantID string `json:"merchant_id"`
		Token      string `json:"token"`
	} `json:"data"`
}

type CreatePaymentRequest struct {
	MerchantID              string  `json:"merchant_id"`
	MerchantKey             string  `json:"merchant_key"`
	MerchantCode            string  `json:"merchant_code"`
	MerchantName            string  `json:"merchant_name"`
	MerchantMobile          string  `json:"merchant_mobile"`
	MerchantCity            string  `json:"merchant_city"`
	MerchantCategoryCode    string  `json:"merchant_category_code"`
	TransactionCurrencyCode string  `json:"transaction_currency_code"`
	Amount                  float64 `json:"amount"`
	InvoiceID               string  `json:"invoice_id"`
	TxnID                   string  `json:"txn_id"`
	Date                    string  `json:"date"`
	RedirectURL             string  `json:"redirect_url"`
	AdditionalInfo          string  `json:"additional_info"`
}

type CreatePaymentResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		GatewayURL string `json:"gateway_url"`
	} `json:"data"`
}

type QueryPaymentRequest struct {
	MerchantID string `json:"merchant_id"`
	InvoiceID  string `json:"invoice_id"`
}

type QueryPaymentResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		InvoiceID       string `json:"invoice_id"`
		TxnID           string `json:"txn_id"`
		UpayTxnID       string `json:"upay_txn_id"`
		Amount          string `json:"amount"`
		Status          string `json:"status"`
		TransactionDate string `json:"transaction_date"`
	} `json:"data"`
}
