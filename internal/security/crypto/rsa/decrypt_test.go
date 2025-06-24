package rsa

import (
	"crypto/rsa"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRsaDecryptor_Decrypt(t *testing.T) {
	file, _ := os.OpenFile("priv_test.pem", os.O_WRONLY|os.O_CREATE, 0766)
	_, _ = file.Write([]byte("testkey"))

	privateKey, _ := LoadPrivateKey("priv_test.pem")

	tests := []struct {
		name       string
		privateKey *rsa.PrivateKey
		ciphertext []byte
		want       []byte
	}{
		{
			name:       "test1",
			privateKey: privateKey,
			ciphertext: []byte("test1"),
			want:       []byte("test1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsaDec := RsaDecryptor{
				PrivateKey: tt.privateKey,
			}
			got, err := rsaDec.Decrypt(tt.ciphertext)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
	_ = os.Remove("priv_test.pem")
}
