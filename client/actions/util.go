package actions

import (
	"crypto/sha512"
	"io/ioutil"
)

func GetSourceChecksum(files []File) [64]byte {
	checksums := make([][64]byte, len(files))
	for i := 0; i < len(files); i++ {
		s, err := ioutil.ReadFile(files[i].GetFullPath())
		if err != nil {
			//panic(err)
			s = make([]byte, 0)
		}
		checksums[i] = sha512.Sum512(s)
	}

	checksum := make([]byte, len(files)*64)
	for i := 0; i < len(files); i++ {
		for j := 0; j < 64; j++ {
			checksum[i*64+j] = checksums[i][j]
		}
	}

	return sha512.Sum512(checksum)
}
