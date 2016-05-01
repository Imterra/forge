package proto

type FileInfo struct {
	Filename string
	Checksum [64]byte
}

type Args struct {
	Name        string
	Inputs      []FileInfo
	SendContent bool
}

type Response struct {
	ActionOutput string
	FileContents []byte
}
