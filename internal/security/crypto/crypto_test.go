package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRSADecryptor(t *testing.T) {
	privateFile, _ := os.OpenFile("priv_test.pem", os.O_WRONLY|os.O_CREATE, 0766)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	bytes, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: bytes,
	}
	_ = pem.Encode(privateFile, privateBlock)

	tests := []struct {
		name           string
		privateKeyPath string
		want           string
	}{
		{
			name:           "test1",
			privateKeyPath: "priv_test.pem",
			want:           "*rsa.RsaDecryptor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRSADecryptor(tt.privateKeyPath)
			assert.Equal(t, tt.want, reflect.ValueOf(got).Type().String())
			assert.NotNil(t, got)
			assert.NoError(t, err)
		})
	}
	_ = os.Remove("priv_test.pem")
}

func TestNewRSAEncryptor(t *testing.T) {
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
		name          string
		publicKeyPath string
		want          string
	}{
		{
			name:          "test1",
			publicKeyPath: "pub_test.pem",
			want:          "*rsa.RsaEncryptor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRSAEncryptor(tt.publicKeyPath)
			assert.Equal(t, tt.want, reflect.ValueOf(got).Type().String())
			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
	_ = os.Remove("pub_test.pem")
}
