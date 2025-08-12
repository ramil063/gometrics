package crypto

import (
	"sync"

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

// Manager содержит все шифровальщики и дешифровщики
type Manager struct {
	defaultEncryptor Encryptor
	defaultDecryptor Decryptor
	grpcEncryptor    Encryptor
	grpcDecryptor    Decryptor
	mx               sync.RWMutex
}

func NewCryptoManager() *Manager {
	return &Manager{}
}

func (cm *Manager) SetDefaultEncryptor(enc Encryptor) {
	cm.mx.RLock()
	defer cm.mx.RUnlock()
	cm.defaultEncryptor = enc
}

func (cm *Manager) SetDefaultDecryptor(decr Decryptor) {
	cm.mx.RLock()
	defer cm.mx.RUnlock()
	cm.defaultDecryptor = decr
}

func (cm *Manager) SetGRPCEncryptor(enc Encryptor) {
	cm.mx.RLock()
	defer cm.mx.RUnlock()
	cm.grpcEncryptor = enc
}

func (cm *Manager) SetGRPCDecryptor(decr Decryptor) {
	cm.mx.RLock()
	defer cm.mx.RUnlock()
	cm.grpcDecryptor = decr
}

func (cm *Manager) GetDefaultEncryptor() Encryptor {
	cm.mx.RLock()
	defer cm.mx.RUnlock()
	return cm.defaultEncryptor
}

func (cm *Manager) GetDefaultDecryptor() Decryptor {
	cm.mx.RLock()
	defer cm.mx.RUnlock()
	return cm.defaultDecryptor
}

func (cm *Manager) GetGRPCEncryptor() Encryptor {
	cm.mx.RLock()
	defer cm.mx.RUnlock()
	return cm.grpcEncryptor
}

func (cm *Manager) GetGRPCDecryptor() Decryptor {
	cm.mx.RLock()
	defer cm.mx.RUnlock()
	return cm.grpcDecryptor
}

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
