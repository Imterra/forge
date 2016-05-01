package actions

import (
	"../../proto"
	"fmt"
	"io/ioutil"
	"strings"
)

func GetInfileData(files []File, rootdir string) ([]proto.FileInfo, error) {
	infiles := make([]proto.FileInfo, len(files))
	for i := range files {
		infiles[i].Filename = files[i].GetPath()
		d, err := ioutil.ReadFile(files[i].GetAbsolutePath(rootdir))
		if err != nil {
			return nil, err
		}
		infiles[i].Checksum = proto.GetDataChecksum(d)
	}
	return infiles, nil
}

func MakeCObjects(name string, sources []string) []File {
	inputs := make([]File, len(sources))

	for i := range sources {
		file := SourceFile{Filename: strings.TrimPrefix(sources[i], "//")}
		action := Action{
			Name:    fmt.Sprintf("CC(%s)", sources[i]),
			Infiles: []File{&file},
			Method:  "Task.CompileC",
		}
		genfile := GeneratedFile{
			Filename: strings.TrimSuffix(
				strings.TrimPrefix(sources[i], "//"), ".c") + ".o",
			Action: &action,
		}
		inputs[i] = &genfile
	}

	return inputs
}
