package actions

import (
	"../../log"
	"../../proto"
	"../util"
	"../worker"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GetInfileData(files []*File, rootdir string) ([]proto.FileInfo, error) {
	infiles := make([]proto.FileInfo, len(files))
	for i := range files {
		infiles[i].Filename = files[i].Filename
		d, err := ioutil.ReadFile(files[i].GetAbsolutePath(rootdir))
		if err != nil {
			return nil, err
		}
		infiles[i].Checksum = proto.GetDataChecksum(d)
	}
	return infiles, nil
}

func MakeCObjects(name string, sources []string, headers []string, file_list map[string]*File) []*File {
	inputs := make([]*File, len(sources))

	header_files := make([]*File, len(headers))
	for i := range headers {
		f, ok := file_list[headers[i]]
		if ok {
			header_files[i] = f
			continue
		}
		file := &File{Filename: headers[i], Action: nil, Sem: make(chan int, 1)}
		file_list[headers[i]] = file
		header_files[i] = file
	}

	for i := range sources {
		filename := strings.TrimPrefix(sources[i], "//")
		outfilename := strings.TrimSuffix(filename, ".c") + ".o"

		infiles := make([]*File, len(header_files)+1)

		f, ok := file_list[outfilename]
		if ok {
			inputs[i] = f
			continue
		}

		f, ok = file_list[filename]
		var file *File
		if ok {
			file = f
		} else {
			file = &File{Filename: filename, Action: nil, Sem: make(chan int, 1)}
			file_list[filename] = file
		}

		infiles[0] = file
		for i := range header_files {
			infiles[i+1] = header_files[i]
		}

		action := Action{
			Name:    fmt.Sprintf("CC(%s)", sources[i]),
			Infiles: infiles,
			Method:  "Task.CompileC",
		}
		genfile := File{
			Filename: outfilename,
			Action:   &action,
			Sem:      make(chan int, 1),
		}
		file_list[outfilename] = &genfile
		inputs[i] = &genfile
	}

	return inputs
}

func SendFile(file *File, worker *worker.Worker, conf *util.Config) error {
	_, ok := worker.Files[file.Filename]

	if ok {
		return nil
	}

	if !worker.Request {
		worker.Files[file.Filename] = 1
		return nil
	}

	var resp proto.FileResponse

	fi, err := os.Stat(file.GetAbsolutePath(conf.Rootdir))
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(file.GetAbsolutePath(conf.Rootdir))
	if err != nil {
		return err
	}

	args := proto.File{
		Filename: file.Filename,
		Content:  content,
		Mode:     fi.Mode(),
	}

	err = worker.Client.Call("File.RecvFile", args, &resp)
	if err != nil {
		return err
	}

	GiveFile(worker, file.Filename)
	return nil
}

var worker_sem chan int = make(chan int, 1)

func ChooseBestWorker(infiles []*File, workers []*worker.Worker) *worker.Worker {
	worker_sem <- 1

	if len(workers) == 0 {
		log.Error("no worker specified!", util.Exiter)
	}

	scores := make([]int, len(workers))

	for iw := range workers {
		w := workers[iw]
		scores[iw] = w.NumTasks

		for i := range infiles {
			_, ok := w.Files[infiles[i].Filename]
			if ok {
				scores[iw]++
			}
		}
	}

	max_i := 0

	for i := range scores {
		if scores[i] > scores[max_i] {
			max_i = i
		}
	}

	chosen := workers[max_i]
	chosen.NumTasks--
	<-worker_sem

	return chosen
}

func FreeWorker(worker *worker.Worker) {
	worker_sem <- 1
	worker.NumTasks++
	<-worker_sem
}

func GiveFile(worker *worker.Worker, filename string) {
	worker_sem <- 1
	worker.Files[filename] = 1
	<-worker_sem
}
