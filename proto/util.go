package proto

import (
	"crypto/sha512"
)

func GetDataChecksum(data []byte) [64]byte {
	return sha512.Sum512(data)
}
