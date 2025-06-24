package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
)

// RsaDecryptor шифровальщик
type RsaDecryptor struct {
	PrivateKey *rsa.PrivateKey
}

// Decrypt функция дешифровки
func (rsaDec RsaDecryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	if rsaDec.PrivateKey == nil {
		return ciphertext, nil
	}

	var result []byte
	pos := 0

	for pos < len(ciphertext) {
		// Проверяем достаточно ли данных для заголовка
		if pos+blockHeaderSize > len(ciphertext) {
			return nil, errors.New("corrupted ciphertext: incomplete block header")
		}

		// Читаем размер зашифрованного блока
		blockSize := int(binary.BigEndian.Uint32(ciphertext[pos : pos+blockHeaderSize]))
		pos += blockHeaderSize

		// Проверяем достаточно ли данных для блока
		if pos+blockSize > len(ciphertext) {
			return nil, errors.New("corrupted ciphertext: incomplete block data")
		}

		// Дешифруем блок
		decrypted, err := rsa.DecryptOAEP(
			sha256.New(),
			rand.Reader,
			rsaDec.PrivateKey,
			ciphertext[pos:pos+blockSize],
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("decryption failed at block %d: %w", len(result)/rsaOAEPSizeLimit, err)
		}

		result = append(result, decrypted...)
		pos += blockSize
	}

	return result, nil
}
