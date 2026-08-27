# Nagad Payment Gateway Integration for Go

A modular and reusable Go library for integrating the Nagad payment gateway. This library handles the complex two-step handshake, RSA/AES encryption, and signature verification required by Nagad.

## Features
- **Two-Step Handshake**: Supports `Initialize` and `Complete` flow.
- **Robust Cryptography**: Built-in support for:
  - RSA Encryption/Decryption (PKCS#1 v1.5)
  - RSA Signing (SHA256withRSA)
  - AES-CBC Encryption (PKCS#7 padding)
- **Environment Support**: Easy switching between Sandbox and Production.
- **Microservice Ready**: Designed as a standalone package for easy integration.

## Project Structure
- `nagad/`: The core library package.
- `examples/`: Usage demonstration.

## Installation
```bash
go get github.com/souravroyscr/nagad-integration
```

## Quick Start

### 1. Configure the Client
```go
config := &nagad.Config{
    MerchantID:         "YOUR_MERCHANT_ID",
    MerchantPrivateKey: "-----BEGIN RSA PRIVATE KEY-----\n...",
    NagadPublicKey:     "-----BEGIN PUBLIC KEY-----\n...",
    APIKey:             "YOUR_API_KEY",
    Environment:        nagad.Sandbox, // or nagad.Production
}
client := nagad.NewClient(config)
```

### 2. Initialize Payment
```go
sensitive, err := client.Initialize("ORDER_ID", "CLIENT_IP")
if err != nil {
    log.Fatal(err)
}
```

### 3. Complete Payment
```go
details := nagad.PaymentDetails{
    MerchantID: config.MerchantID,
    OrderID:    "ORDER_ID",
    Amount:     "100.00",
    Currency:   "BDT",
    Challenge:  "random_string",
}

resp, err := client.Complete(sensitive.PaymentReferenceID, details, *sensitive, "CLIENT_IP", "CALLBACK_URL")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Redirect user to: %s\n", resp.CallBackURL)
```

### 4. Verify Payment (IPN or Callback)
```go
verifyResp, err := client.Verify("PAYMENT_REF_ID")
if err != nil {
    log.Fatal(err)
}
if verifyResp.Status == "Success" {
    // Fulfill order
}
```

## Security Note
- Always keep your `MerchantPrivateKey` and `APIKey` secret.
- Use environment variables or a secure vault to manage credentials.
- Do not log sensitive data in production.

## Testing
Run unit tests for cryptographic utilities:
```bash
go test ./nagad/...
```
