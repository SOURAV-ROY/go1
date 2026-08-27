package nagad

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// EncryptRSA encrypts data using a public key in PKCS#1 v1.5 format
func EncryptRSA(publicKeyPEM string, data []byte) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not an RSA public key")
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, data)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt: %w", err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptRSA decrypts data using a private key in PKCS#1 v1.5 format
func DecryptRSA(privateKeyPEM string, ciphertextBase64 string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try parsing as PKCS#8 if PKCS#1 fails
		key, err8 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err8 != nil {
			return nil, fmt.Errorf("failed to parse private key (tried PKCS#1 and PKCS#8): %v, %v", err, err8)
		}
		var ok bool
		priv, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not an RSA private key")
		}
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// SignRSA generates a SHA256withRSA signature
func SignRSA(privateKeyPEM string, data []byte) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		key, err8 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err8 != nil {
			return "", fmt.Errorf("failed to parse private key: %w", err8)
		}
		var ok bool
		priv, ok = key.(*rsa.PrivateKey)
		if !ok {
			return "", errors.New("not an RSA private key")
		}
	}

	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed)
	if err != nil {
		return "", fmt.Errorf("failed to sign: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// AESEncrypt encrypts data using AES-CBC with PKCS#7 padding
func AESEncrypt(key, iv, data []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	content := pkcs7Padding(data, block.BlockSize())
	ciphertext := make([]byte, len(content))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, content)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
