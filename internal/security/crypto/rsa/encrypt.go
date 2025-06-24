package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// RsaEncryptor шифровальщик
type RsaEncryptor struct {
	PublicKey *rsa.PublicKey
}

// Encrypt функция шифрования
func (rsaEnc RsaEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	if rsaEnc.PublicKey == nil {
		return plaintext, nil
	}

	var result []byte

	// Разбиваем данные на блоки
	for offset := 0; offset < len(plaintext); offset += rsaOAEPSizeLimit {
		end := offset + rsaOAEPSizeLimit
		if end > len(plaintext) {
			end = len(plaintext)
		}
		block := plaintext[offset:end]

		encrypted, err := rsa.EncryptOAEP(
			sha256.New(),
			rand.Reader,
			rsaEnc.PublicKey,
			block,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("encryption failed at block %d: %w", offset/rsaOAEPSizeLimit, err)
		}

		// Добавляем размер блока (4 байта) и данные
		sizeBuf := make([]byte, blockHeaderSize)
		binary.BigEndian.PutUint32(sizeBuf, uint32(len(encrypted)))
		result = append(result, sizeBuf...)
		result = append(result, encrypted...)
	}

	return result, nil
}
