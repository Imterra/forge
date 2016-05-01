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

func MakeCObjects(name string, sources []string, file_list *map[string]File) []File {
	inputs := make([]File, len(sources))

	for i := range sources {
		filename := strings.TrimPrefix(sources[i], "//")
		outfilename := strings.TrimSuffix(filename, ".c") + ".o"

		f, ok := (*file_list)[outfilename]
		if ok {
			inputs[i] = f
			continue
		}

		f, ok = (*file_list)[filename]
		var file File
		if ok {
			file = f
		} else {
			file = &SourceFile{Filename: filename}
			(*file_list)[filename] = file
		}

		action := Action{
			Name:    fmt.Sprintf("CC(%s)", sources[i]),
			Infiles: []File{file},
			Method:  "Task.CompileC",
		}
		genfile := GeneratedFile{
			Filename: outfilename,
			Action:   &action,
		}
		(*file_list)[outfilename] = &genfile
		inputs[i] = &genfile
	}

	return inputs
}
