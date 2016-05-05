package actions

import (
	"../../proto"
	"../util"
	"../worker"
	"log"
	"net/rpc"
)

type File struct {
	Filename string
	Action   *Action
	Sem      chan int
}

func (f *File) GetAbsolutePath(rootdir string) string {
	return util.GetFullPath(f.Filename, rootdir)
}

type Action struct {
	Name    string
	Infiles []*File
	Method  string
}

func (action *Action) Execute(client *rpc.Client, config *util.Config) *rpc.Call {
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

	return client.Go(action.Method, &args, &resp, nil)
}

func MakeFile(file *File, conf *util.Config, notify chan *File) {
	file.Sem <- 1

	action := file.Action
	file.Action = nil

	if action == nil {
		// TODO: Check if file exists, notify.
		<-file.Sem
		notify <- file
		return
	}

	new_notify := make(chan *File, len(action.Infiles))

	for i := range action.Infiles {
		go MakeFile(action.Infiles[i], conf, new_notify)
	}

	// TODO: Chose best worker.
	var worker worker.Worker

	for _ = range action.Infiles {
		// TODO: Check if file exists, send to worker.
	}

	call := action.Execute(worker.Client, conf)
	// TODO: Wait for action to finish, check if successful, if so, notify, otherwise, log error.
	<-call.Done

	<-file.Sem
	notify <- file
}
