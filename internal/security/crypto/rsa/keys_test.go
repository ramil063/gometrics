package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPrivateKey(t *testing.T) {
	privateFile, _ := os.Create("priv_test.pem")
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	bytes, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: bytes,
	}
	_ = pem.Encode(privateFile, privateBlock)

	tests := []struct {
		name string
		path string
	}{
		{
			name: "test1",
			path: "priv_test.pem",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadPrivateKey(tt.path)
			assert.NoError(t, err)
			assert.NotEmpty(t, got.PublicKey)
		})
	}
	_ = os.Remove("priv_test.pem")
}

func TestLoadPublicKey(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	publicKey := &privateKey.PublicKey
	publicBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	publicFile, _ := os.Create("pub_test.pem")
	defer publicFile.Close()

	publicBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicBytes,
	}

	_ = pem.Encode(publicFile, publicBlock)

	tests := []struct {
		name string
		path string
	}{
		{
			name: "test1",
			path: "pub_test.pem",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadPublicKey(tt.path)
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
	_ = os.Remove("pub_test.pem")
}
