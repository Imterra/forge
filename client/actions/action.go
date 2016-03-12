package actions

type File interface {
	GetName() string
	GetFullPath() string
	GetOrigin() Action
}

type SourceFile struct {
	Name     string
	FullPath string
}

func (f *SourceFile) GetName() string {
	return f.Name
}

func (f *SourceFile) GetFullPath() string {
	return f.FullPath
}

func (f *SourceFile) GetOrigin() Action {
	return nil
}

type GeneratedFile struct {
	Name     string
	FullPath string
	Origin   Action
}

func (f *GeneratedFile) GetName() string {
	return f.Name
}

func (f *GeneratedFile) GetFullPath() string {
	return f.FullPath
}

func (f *GeneratedFile) GetOrigin() Action {
	return f.Origin
}

type Action interface {
	GetName() string
	GetInfiles() []File
	GetCmd() string
	GetOutFileName() string
}
