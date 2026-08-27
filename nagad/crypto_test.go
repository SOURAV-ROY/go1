package nagad

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestCrypto(t *testing.T) {
	// 1. Generate a temporary RSA key pair for testing
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	// 2. Test RSA Encryption/Decryption
	plaintext := []byte("Hello Nagad!")
	ciphertextBase64, err := EncryptRSA(string(pubKeyPEM), plaintext)
	if err != nil {
		t.Fatalf("RSA Encryption failed: %v", err)
	}

	decrypted, err := DecryptRSA(string(privKeyPEM), ciphertextBase64)
	if err != nil {
		t.Fatalf("RSA Decryption failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("RSA Decryption mismatch. Got %s, want %s", string(decrypted), string(plaintext))
	}

	// 3. Test RSA Signing
	signature, err := SignRSA(string(privKeyPEM), plaintext)
	if err != nil {
		t.Fatalf("RSA Signing failed: %v", err)
	}
	if signature == "" {
		t.Error("RSA Signature is empty")
	}

	// 4. Test AES Encryption
	key := []byte("1234567890123456") // 16 bytes
	iv := []byte("1234567890123456")  // 16 bytes
	aesCiphertext, err := AESEncrypt(key, iv, plaintext)
	if err != nil {
		t.Fatalf("AES Encryption failed: %v", err)
	}
	if aesCiphertext == "" {
		t.Error("AES Ciphertext is empty")
	}
}
