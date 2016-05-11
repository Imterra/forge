package worker

import (
	"net/rpc"
)

type Worker struct {
	Client   *rpc.Client
	Addr     string
	NumTasks int
	Files    map[string]int
	Request  bool
}

func GetWorker(addr string) (*Worker, error) {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	args := 0
	var num_tasks int

	err = client.Call("Util.NumTasks", args, &num_tasks)
	if err != nil {
		return nil, err
	}

	w := Worker{
		Client:   client,
		Addr:     addr,
		NumTasks: num_tasks,
		Files:    make(map[string]int),
		Request:  true,
	}
	return &w, nil
}
