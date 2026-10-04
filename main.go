package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/souravroyscr/nagad-integration/nagad"
	"github.com/souravroyscr/nagad-integration/upay"
)

var (
	nagadClient *nagad.Client
	nagadConfig *nagad.Config
	upayClient  *upay.Client
	upayConfig  *upay.Config
)

func init() {
	// Initialize Nagad
	nagadConfig = &nagad.Config{
		MerchantID:         getEnv("NAGAD_MERCHANT_ID", "683002007104225"),
		MerchantPrivateKey: getEnv("NAGAD_MERCHANT_PRIVATE_KEY", "-----BEGIN RSA PRIVATE KEY-----\n..."),
		NagadPublicKey:     getEnv("NAGAD_PUBLIC_KEY", "-----BEGIN PUBLIC KEY-----\n..."),
		APIKey:             getEnv("NAGAD_API_KEY", "8515453154515451"),
		Environment:        nagad.Sandbox,
	}
	nagadClient = nagad.NewClient(nagadConfig)

	// Initialize UPAY
	upayConfig = &upay.Config{
		MerchantID:           getEnv("UPAY_MERCHANT_ID", "4578753234567"),
		MerchantKey:          getEnv("UPAY_MERCHANT_KEY", "ADSE1234"),
		MerchantCode:         getEnv("UPAY_MERCHANT_CODE", "1252"),
		MerchantName:         getEnv("UPAY_MERCHANT_NAME", "Test Store"),
		MerchantMobile:       getEnv("UPAY_MERCHANT_MOBILE", "017XXXXXXXX"),
		MerchantCity:         getEnv("UPAY_MERCHANT_CITY", "Dhaka"),
		MerchantCategoryCode: getEnv("UPAY_MERCHANT_CATEGORY_CODE", "1252"),
		Environment:          upay.Sandbox,
	}
	upayClient = upay.NewClient(upayConfig)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/pay", handlePay)
	mux.HandleFunc("/api/verify", handleVerify)

	// Add middleware
	handler := loggingMiddleware(corsMiddleware(mux))

	port := getEnv("PORT", "8080")
	fmt.Printf("Bridge server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func handlePay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Gateway string `json:"gateway"`
		OrderID string `json:"orderId"`
		Amount  string `json:"amount"`
		IP      string `json:"ip"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Gateway == "upay" {
		amountFloat := 0.0
		fmt.Sscanf(req.Amount, "%f", &amountFloat)
		callbackURL := getEnv("UPAY_CALLBACK_URL", "http://localhost:5173/callback?gateway=upay")
		resp, err := upayClient.CreatePayment(req.OrderID, amountFloat, callbackURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("UPAY Initiation failed: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"callBackUrl": resp.Data.GatewayURL,
			"status":      "Success",
		})
		return
	}

	// Default to Nagad
	// 1. Initialize
	sensitive, err := nagadClient.Initialize(req.OrderID, req.IP)
	if err != nil {
		http.Error(w, fmt.Sprintf("Initialization failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 2. Complete
	details := nagad.PaymentDetails{
		MerchantID: nagadConfig.MerchantID,
		OrderID:    req.OrderID,
		Amount:     req.Amount,
		Currency:   "BDT",
		Challenge:  "random_challenge_string",
	}

	callbackURL := getEnv("NAGAD_CALLBACK_URL", "http://localhost:5173/callback?gateway=nagad")
	resp, err := nagadClient.Complete(sensitive.PaymentReferenceID, details, *sensitive, req.IP, callbackURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Completion failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Gateway      string `json:"gateway"`
		PaymentRefID string `json:"paymentRefId"` // For Nagad
		InvoiceID    string `json:"invoiceId"`    // For UPAY
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Gateway == "upay" {
		resp, err := upayClient.QueryPayment(req.InvoiceID)
		if err != nil {
			http.Error(w, fmt.Sprintf("UPAY Verification failed: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": resp.Data.Status,
			"data":   resp.Data,
		})
		return
	}

	// Default to Nagad
	resp, err := nagadClient.Verify(req.PaymentRefID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Verification failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// responseWriter is a wrapper for http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		log.Printf(
			"method=%s path=%s remote=%s status=%d duration=%s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			rw.status,
			time.Since(start),
		)
	})
}
