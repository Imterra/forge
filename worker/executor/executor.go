package executor

import (
	"fmt"
)

type Config struct {
	Pkgroot string
	Outdir  string
}

//type Task struct {
//	Type        string
//	OutFilename string
//	SrcArgs     []string
//	ResArgs     []string
//	DepArgs     []string
//}

type Task interface {
	GetType() string
	GetOutfile() string
	GetInputs() []string
	Execute(config *Config) (string, error)
}

type QuitTask struct{}

func (t *QuitTask) GetType() string {
	return "stop"
}
func (t *QuitTask) GetOutfile() string {
	return ""
}
func (t *QuitTask) GetInputs() []string {
	return make([]string, 0)
}
func (t *QuitTask) Execute(config *Config) (string, error) {
	return "", nil
}

const (
	WILCO = iota
	DONE
	QUIT
	ERROR
)

type Response struct {
	Id      int
	Code    int
	Output  string
	Outfile string
}

func Executor(id int, taskqueue chan Task, responder chan Response, config *Config) {
	var task Task

	//WILCO_RESPONSE := Response{Id: id, Code: WILCO, Output: "", Outfile: ""}
	QUIT_RESPONSE := Response{Id: id, Code: QUIT, Output: "", Outfile: ""}

	for {
		task = <-taskqueue
		if task.GetType() == "stop" {
			responder <- QUIT_RESPONSE
			return
		}
		output, err := task.Execute(config)
		if err != nil {
			responder <- Response{Id: id, Code: ERROR, Output: fmt.Sprintf("%v\n%v\n", err, output), Outfile: task.GetOutfile()}
		}
		responder <- Response{Id: id, Code: DONE, Output: output, Outfile: task.GetOutfile()}
	}
}
