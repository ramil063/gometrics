package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// CreateSha256 создаем хеш
func CreateSha256(body []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)
}
