package actions

import (
	"../../log"
	"../../proto"
	"../util"
	"../worker"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func (action *Action) Execute(config *util.Config, worker *worker.Worker) error {
	infiles, err := GetInfileData(action.Infiles, config.Rootdir)

	if err != nil {
		log.Error(fmt.Sprintf("executing action %s failed: %s", action.Name, err.Error()))
	}

	args := proto.Args{
		Name:        action.Name,
		Inputs:      infiles,
		SendContent: worker.Request,
	}
	var resp = new(proto.Response)

	return worker.Client.Call(action.Method, &args, &resp)
}

func MakeFile(file *File, conf *util.Config, notify chan *File) {
	file.Sem <- 1
	fmt.Printf("[DBG] working on file: %s\n", file.Filename)
	defer func() {
		<-file.Sem
		notify <- file
	}()

	action := file.Action
	file.Action = nil

	if action == nil {
		_, err := os.Stat(file.GetAbsolutePath(conf.Rootdir))
		if err != nil {
			log.Error(fmt.Sprintf("file %s does not exist: %s", file.Filename, err.Error()))
		}
		return
	}

	new_notify := make(chan *File, len(action.Infiles))

	for i := range action.Infiles {
		fmt.Printf("[DBG] requesting file: %s\n", action.Infiles[i].Filename)
		go func(i int) {
			defer log.HandleExit()
			MakeFile(action.Infiles[i], conf, new_notify)
		}(i)
	}

	rebuild := false

	for _ = range action.Infiles {
		f := <-new_notify

		fmt.Printf("[DBG] got file: %s for action: %s\n", f.Filename, action.Name)

		_, err := os.Stat(f.GetAbsolutePath(conf.Rootdir))
		if err != nil {
			log.Error(fmt.Sprintf("file %s does not exist: %s", f.Filename, err.Error()))
		}

		checksum, err := GetFileChecksum(f, conf)
		if err != nil {
			rebuild = true
			continue
		}

		meta_checksum, err := GetMetadata(f, conf)
		if err != nil {
			rebuild = true
			continue
		}

		if checksum != meta_checksum {
			rebuild = true
			continue
		}
	}

	checksum, err := GetFileChecksum(file, conf)
	if err != nil {
		rebuild = true
	}

	meta_checksum, err := GetMetadata(file, conf)
	if err != nil {
		rebuild = true
	}

	if checksum != meta_checksum {
		rebuild = true
	}

	if !rebuild {
		return
	}

	worker := ChooseBestWorker(action.Infiles, conf.Workers)

	for i := range action.Infiles {
		f := action.Infiles[i]
		err := SendFile(f, worker, conf)
		if err != nil {
			log.Error(fmt.Sprintf("error sending file %s: %s", err.Error()))
		}
	}

	err = action.Execute(conf, worker)
	if err != nil {
		log.Error(
			fmt.Sprintf(
				"executing action for file %s: %s", file.Filename, err.Error()))
	}
	GiveFile(worker, file.Filename)
	FreeWorker(worker)

	if worker.Request {

		request := proto.FileRequest{Filename: file.Filename}
		var resp proto.File

		err = worker.Client.Call("File.SendFile", request, &resp)
		if err != nil {
			log.Error(fmt.Sprintf("receiving file: %s: %s", file.Filename, err.Error()))
		}

		full_path := file.GetAbsolutePath(conf.Rootdir)
		full_dir := filepath.Dir(full_path)

		var mode os.FileMode = os.ModeDir + 0755
		err = os.MkdirAll(full_dir, mode)
		if err != nil {
			log.Error(fmt.Sprintf("receiving file: %s: %s", file.Filename, err.Error()))
		}

		err = ioutil.WriteFile(full_path, resp.Content, resp.Mode)
		if err != nil {
			log.Error(fmt.Sprintf("receiving file: %s: %s", file.Filename, err.Error()))
		}

	}
	log.Succ(file.Filename)
}

func GetFileChecksum(file *File, conf *util.Config) (string, error) {
	path := file.GetAbsolutePath(conf.Rootdir)
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	chksum := proto.GetDataChecksum(d)
	s := ""
	for i := range chksum {
		s = fmt.Sprintf("%s%02x", s, chksum[i])
	}
	return s, nil
}

func GetMetadata(file *File, config *util.Config) (string, error) {
	metadir := filepath.Join(config.Rootdir, ".metadata")
	metafilename := file.GetAbsolutePath(metadir)

	d, err := ioutil.ReadFile(metafilename)
	if err != nil {
		return "", err
	}

	return string(d), nil
}
