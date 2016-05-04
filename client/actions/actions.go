package actions

import (
	"../../proto"
	"../util"
	"log"
	"net/rpc"
)

type File struct {
	Filename string
	Action   *Action
}

func (f *File) GetAbsolutePath(rootdir string) string {
	return util.GetFullPath(f.Filename, rootdir)
}

type Action struct {
	Name    string
	Infiles []*File
	Method  string
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
