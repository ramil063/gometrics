package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"strings"
	"testing"
)

const (
	blockHeaderSizeTest = 4
)

// Генерация тестового RSA ключа
func generateTestKey(t *testing.T) *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}
	return key
}

// Создание зашифрованного блока для тестов
func createTestBlock(t *testing.T, key *rsa.PublicKey, data []byte) []byte {
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, data, nil)
	if err != nil {
		t.Fatalf("Failed to create test block: %v", err)
	}

	header := make([]byte, blockHeaderSizeTest)
	binary.BigEndian.PutUint32(header, uint32(len(encrypted)))

	return append(header, encrypted...)
}

func TestRsaDecryptor_Decrypt(t *testing.T) {
	key := generateTestKey(t)
	testData := []byte("test data to encrypt")

	tests := []struct {
		name        string
		setup       func() (RsaDecryptor, []byte)
		want        []byte
		wantErr     bool
		errContains string
	}{
		{
			name: "nil private key returns original",
			setup: func() (RsaDecryptor, []byte) {
				return RsaDecryptor{PrivateKey: nil}, []byte("test")
			},
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name: "successful single block decryption",
			setup: func() (RsaDecryptor, []byte) {
				block := createTestBlock(t, &key.PublicKey, testData)
				return RsaDecryptor{PrivateKey: key}, block
			},
			want:    testData,
			wantErr: false,
		},
		{
			name: "corrupted header",
			setup: func() (RsaDecryptor, []byte) {
				return RsaDecryptor{PrivateKey: key}, []byte{0, 0} // Неполный заголовок
			},
			wantErr:     true,
			errContains: "incomplete block header",
		},
		{
			name: "corrupted block data",
			setup: func() (RsaDecryptor, []byte) {
				header := make([]byte, blockHeaderSize)
				binary.BigEndian.PutUint32(header, 100) // Указываем большой размер блока
				return RsaDecryptor{PrivateKey: key}, header
			},
			wantErr:     true,
			errContains: "incomplete block data",
		},
		{
			name: "multi-block decryption",
			setup: func() (RsaDecryptor, []byte) {
				// Создаем 2 блока
				block1 := createTestBlock(t, &key.PublicKey, []byte("first part"))
				block2 := createTestBlock(t, &key.PublicKey, []byte("second part"))
				return RsaDecryptor{PrivateKey: key}, append(block1, block2...)
			},
			want:    []byte("first partsecond part"),
			wantErr: false,
		},
		{
			name: "decryption failure",
			setup: func() (RsaDecryptor, []byte) {
				// Создаем валидный блок, но с другим ключом
				wrongKey := generateTestKey(t)
				block := createTestBlock(t, &wrongKey.PublicKey, testData)
				return RsaDecryptor{PrivateKey: key}, block
			},
			wantErr:     true,
			errContains: "decryption failed at block 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decryptor, input := tt.setup()
			got, err := decryptor.Decrypt(input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Decrypt() error = %v, want contains %v", err, tt.errContains)
				}
				return
			}

			if !bytes.Equal(got, tt.want) {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
