package proto

import (
	"os"
)

type File struct {
	Filename string
	Content  []byte
	Mode     os.FileMode
}

type FileRequest struct {
	Filename string
}

type FileResponse struct {
	Filename    string
	Checksum    [64]byte
	Overwritten bool
}
