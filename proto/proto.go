package proto

type Args struct {
	Name        string
	Inputs      []string
	SendContent bool
}

type Response struct {
	ActionOutput string
	FileContents []byte
}
