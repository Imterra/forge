package actions

import (
	"../../proto"
	"../util"
	"log"
	"net/rpc"
)

type File interface {
	GetPath() string
	GetAbsolutePath(rootdir string) string
	GetOrigin() *Action
}

type Action struct {
	Name    string
	Infiles []File
	Method  string
}

type SourceFile struct {
	Filename string
}

func (f *SourceFile) GetPath() string {
	return f.Filename
}

func (f *SourceFile) GetAbsolutePath(rootdir string) string {
	return util.GetFullPath(f.Filename, rootdir)
}

func (f *SourceFile) GetOrigin() *Action {
	return nil
}

type GeneratedFile struct {
	Filename string
	Action   *Action
}

func (f *GeneratedFile) GetPath() string {
	return f.Filename
}

func (f *GeneratedFile) GetAbsolutePath(rootdir string) string {
	return util.GetFullPath(f.Filename, rootdir)
}

func (f *GeneratedFile) GetOrigin() *Action {
	return f.Action
}

func Execute(action Action, client *rpc.Client, config *util.Config) {
	infiles, err := GetInfileData(action.Infiles, config.Rootdir)

	if err != nil {
		log.Fatal(err)
	}

	args := proto.Args{
		Name:        action.Name,
		Inputs:      infiles,
		SendContent: config.Request,
	}
	var resp = new(proto.Response)

	config.Client.Go(action.Method, &args, &resp, config.Monitor)
}
