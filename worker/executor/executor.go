package executor

import (
	"fmt"
	"os/exec"
	"strings"
)

type Task struct {
	Type        string
	OutFilename string
	SrcArgs     []string
	ResArgs     []string
	DepArgs     []string
}

const (
	WILCO = iota
	DONE
	QUIT
)

type Response struct {
	Id      int
	Code    int
	Outputs []string
}

func Executor(id int, taskqueue chan Task, responder chan Response) {
	var task Task

	WILCO_RESPONSE := Response{Id: id, Code: WILCO, Outputs: make([]string, 0)}
	QUIT_RESPONSE := Response{Id: id, Code: QUIT, Outputs: make([]string, 0)}

	for {
		task = <-taskqueue
		if task.Type == "stop" {
			responder <- QUIT_RESPONSE
			return
		}
		responder <- WILCO_RESPONSE
		responder <- *process_task(id, &task)
	}
}

func process_task(id int, task *Task) *Response {
	switch task.Type {
	case "compile_c":
		s := fmt.Sprintf("gcc -o %s -c %s", task.OutFilename, strings.Join(task.SrcArgs, " "))
		output, err := exec.Command("echo", s).Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("CMD OUT: %s\n", output)
	}
	arr := make([]string, 1)
	arr[0] = task.OutFilename
	return &Response{Id: id, Code: DONE, Outputs: arr}
}
