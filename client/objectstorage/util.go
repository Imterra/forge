package objectstorage

import (
	"fmt"
)

func GetHashString(hash [64]byte) string {
	s := ""
	for i := 0; i < len(hash); i++ {
		s = fmt.Sprintf("%s%02x", s, hash[i])
	}
	return s
}
