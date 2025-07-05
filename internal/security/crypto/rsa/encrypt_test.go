package rsa

import (
	"crypto/rsa"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRsaEncryptor_Encrypt(t *testing.T) {
	file, _ := os.OpenFile("pub_test.pub", os.O_WRONLY|os.O_CREATE, 0766)
	_, _ = file.Write([]byte("testkey"))

	publicKey, _ := LoadPublicKey("pub_test.pub")
	tests := []struct {
		name      string
		publicKey *rsa.PublicKey
		plaintext []byte
		want      []byte
	}{
		{
			name:      "test1",
			publicKey: publicKey,
			plaintext: []byte("test"),
			want:      []byte("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsaEnc := RsaEncryptor{
				PublicKey: tt.publicKey,
			}
			got, err := rsaEnc.Encrypt(tt.plaintext)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
	_ = os.Remove("pub_test.pub")
}
