package tasks

import (
	"../../proto"
	"io/ioutil"
	"path/filepath"
)

func processInputs(inputs []string, dir string, prepend []string) []string {
	ret := make([]string, len(inputs)+len(prepend))
	P := len(prepend)

	for i := range prepend {
		ret[i] = prepend[i]
	}

	for i := range inputs {
		ret[i+P] = filepath.Join(dir, inputs[i])
	}

	return ret
}

func prepareResponse(outfilepath string, output []byte,
	send bool, resp *proto.Response) error {

	data, err := ioutil.ReadFile(outfilepath)
	if err != nil {
		return err
	}

	resp.ActionOutput = string(output)
	if send {
		resp.FileContents = data
	} else {
		resp.FileContents = make([]byte, 0)
	}

	return nil
}
