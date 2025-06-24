package rsa

const (
	// Для RSA 2048 бит с OAEP и SHA-256:
	// Максимальный размер блока = 256 - 2*32 - 2 = 190 байт
	rsaOAEPSizeLimit = 190
	// Размер заголовка блока (4 байта для uint32)
	blockHeaderSize = 4
)
