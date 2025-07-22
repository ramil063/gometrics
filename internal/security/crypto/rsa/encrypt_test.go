package rsa

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"strings"
	"testing"
)

func TestRsaEncryptor_Encrypt(t *testing.T) {
	key := generateTestKey(t)

	tests := []struct {
		name        string
		setup       func() (RsaEncryptor, []byte)
		wantErr     bool
		errContains string
		validate    func(t *testing.T, result []byte)
	}{
		{
			name: "nil public key returns original",
			setup: func() (RsaEncryptor, []byte) {
				return RsaEncryptor{PublicKey: nil}, []byte("test")
			},
			validate: func(t *testing.T, result []byte) {
				if !bytes.Equal(result, []byte("test")) {
					t.Errorf("Expected original plaintext, got %v", result)
				}
			},
		},
		{
			name: "successful single block encryption",
			setup: func() (RsaEncryptor, []byte) {
				return RsaEncryptor{PublicKey: &key.PublicKey}, []byte("short message")
			},
			validate: func(t *testing.T, result []byte) {
				if len(result) <= blockHeaderSize {
					t.Error("Expected encrypted data with header")
				}
			},
		},
		{
			name: "successful multi-block encryption",
			setup: func() (RsaEncryptor, []byte) {
				// Создаем данные больше одного блока
				data := make([]byte, rsaOAEPSizeLimit*2)
				rand.Read(data)
				return RsaEncryptor{PublicKey: &key.PublicKey}, data
			},
			validate: func(t *testing.T, result []byte) {
				// Проверяем что есть несколько блоков
				pos := 0
				blocks := 0
				for pos < len(result) {
					if pos+blockHeaderSize > len(result) {
						t.Error("Incomplete block header")
						return
					}
					blockSize := int(binary.BigEndian.Uint32(result[pos : pos+blockHeaderSize]))
					pos += blockHeaderSize
					if pos+blockSize > len(result) {
						t.Error("Incomplete block data")
						return
					}
					pos += blockSize
					blocks++
				}
				if blocks < 2 {
					t.Errorf("Expected at least 2 blocks, got %d", blocks)
				}
			},
		},
		{
			name: "empty input",
			setup: func() (RsaEncryptor, []byte) {
				return RsaEncryptor{PublicKey: &key.PublicKey}, []byte{}
			},
			validate: func(t *testing.T, result []byte) {
				if len(result) != 0 {
					t.Error("Expected empty result for empty input")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptor, input := tt.setup()
			result, err := encryptor.Encrypt(input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Encrypt() error = %v, want contains %v", err, tt.errContains)
				}
				return
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
