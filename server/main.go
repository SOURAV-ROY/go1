package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/souravroyscr/nagad-integration/nagad"
)

var client *nagad.Client
var config *nagad.Config

func init() {
	// Initialize Nagad Config with environment variables or defaults
	config = &nagad.Config{
		MerchantID:         getEnv("NAGAD_MERCHANT_ID", "683002007104225"),
		MerchantPrivateKey: getEnv("NAGAD_MERCHANT_PRIVATE_KEY", "-----BEGIN RSA PRIVATE KEY-----\n..."),
		NagadPublicKey:     getEnv("NAGAD_PUBLIC_KEY", "-----BEGIN PUBLIC KEY-----\n..."),
		APIKey:             getEnv("NAGAD_API_KEY", "8515453154515451"),
		Environment:        nagad.Sandbox,
	}
	client = nagad.NewClient(config)
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

	// Add basic CORS middleware
	handler := corsMiddleware(mux)

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
		OrderID string `json:"orderId"`
		Amount  string `json:"amount"`
		IP      string `json:"ip"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 1. Initialize
	sensitive, err := client.Initialize(req.OrderID, req.IP)
	if err != nil {
		http.Error(w, fmt.Sprintf("Initialization failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 2. Complete
	details := nagad.PaymentDetails{
		MerchantID: config.MerchantID,
		OrderID:    req.OrderID,
		Amount:     req.Amount,
		Currency:   "BDT",
		Challenge:  "random_challenge_string",
	}

	callbackURL := getEnv("NAGAD_CALLBACK_URL", "http://localhost:5173/callback")
	resp, err := client.Complete(sensitive.PaymentReferenceID, details, *sensitive, req.IP, callbackURL)
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
		PaymentRefID string `json:"paymentRefId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := client.Verify(req.PaymentRefID)
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
