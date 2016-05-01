package actions

import (
	"fmt"
	"strings"
)

func GetInfilePaths(files []File) []string {
	infiles := make([]string, len(files))
	for i := range files {
		infiles[i] = files[i].GetPath()
	}
	return infiles
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
		fmt.Printf("[DBG] TN: %s, IF: %s, OF: %s\n", name, file.Filename, genfile.Filename)
		inputs[i] = &genfile
	}

	return inputs
}
