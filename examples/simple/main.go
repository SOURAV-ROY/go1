package main

import (
	"fmt"
	"log"

	"github.com/souravroyscr/nagad-integration/nagad"
)

func main() {
	// Configuration for Sandbox
	config := &nagad.Config{
		MerchantID:         "683002007104225",
		MerchantPrivateKey: `-----BEGIN RSA PRIVATE KEY-----
... your private key ...
-----END RSA PRIVATE KEY-----`,
		NagadPublicKey: `-----BEGIN PUBLIC KEY-----
... nagad public key ...
-----END PUBLIC KEY-----`,
		APIKey:      "8515453154515451",
		Environment: nagad.Sandbox,
	}

	client := nagad.NewClient(config)

	// Step 1: Initialize
	orderID := "ORDER_12345"
	clientIP := "127.0.0.1"
	
	fmt.Printf("Initializing payment for Order: %s\n", orderID)
	// Note: This will fail without real credentials/connectivity
	sensitive, err := client.Initialize(orderID, clientIP)
	if err != nil {
		log.Printf("Initialization failed (expected with dummy keys): %v\n", err)
		return
	}

	fmt.Printf("Payment Reference ID: %s\n", sensitive.PaymentReferenceID)

	// Step 2: Complete
	details := nagad.PaymentDetails{
		MerchantID: config.MerchantID,
		OrderID:    orderID,
		Amount:     "100.00",
		Currency:   "BDT",
		Challenge:  "random_challenge_string",
	}

	callbackURL := "https://your-microservice.com/nagad/callback"
	resp, err := client.Complete(sensitive.PaymentReferenceID, details, *sensitive, clientIP, callbackURL)
	if err != nil {
		log.Fatalf("Completion failed: %v", err)
	}

	fmt.Printf("Redirect to: %s\n", resp.CallBackURL)
}
