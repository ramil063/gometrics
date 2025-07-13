package crypto

import (
	"github.com/ramil063/gometrics/internal/security/crypto/rsa"
)

// Encryptor общий интерфейс для шифрования
type Encryptor interface {
	Encrypt(plaintext []byte) ([]byte, error)
}

// Decryptor общий интерфейс для шифрования
type Decryptor interface {
	Decrypt(encrypted []byte) ([]byte, error)
}

// DefaultEncryptor стандартный шифровщик
var DefaultEncryptor Encryptor
var DefaultDecryptor Decryptor

// NewRSAEncryptor фабрика для RSA
func NewRSAEncryptor(publicKeyPath string) (Encryptor, error) {
	var encryptor rsa.RsaEncryptor
	var err error

	encryptor.PublicKey, err = rsa.LoadPublicKey(publicKeyPath)
	return &encryptor, err
}

// NewRSADecryptor фабрика для RSA
func NewRSADecryptor(privateKeyPath string) (Decryptor, error) {
	var decryptor rsa.RsaDecryptor
	var err error

	decryptor.PrivateKey, err = rsa.LoadPrivateKey(privateKeyPath)
	return &decryptor, err
}
