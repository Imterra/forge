package target

import (
	"../actions"
	"strings"
)

type Target interface {
	GetName() string
	GetSources() []string
	GetResources() []string
	GetDependencies() []Target
	GetOutputFile() *actions.File
}

var file_list map[string]*actions.File

type LibCTarget struct {
	Name         string
	Sources      []string
	Resources    []string
	Dependencies []Target
}

type AppCTarget struct {
	Name         string
	Sources      []string
	Resources    []string
	Dependencies []Target
}

func (t *LibCTarget) GetName() string {
	return t.Name
}

func (t *LibCTarget) GetSources() []string {
	return t.Sources
}

func (t *LibCTarget) GetResources() []string {
	return t.Resources
}

func (t *LibCTarget) GetDependencies() []Target {
	return t.Dependencies
}

func (t *LibCTarget) GetOutputFile() *actions.File {
	if file_list == nil {
		file_list = make(map[string]*actions.File)
	}

	outfile_name := strings.TrimPrefix(t.Name, "//") + ".a"

	f, ok := file_list[outfile_name]
	if ok {
		return f
	}

	inputs := actions.MakeCObjects(t.Name, t.Sources, file_list)
	ar_action := actions.Action{
		Name:    strings.TrimPrefix(t.Name, "//"),
		Infiles: inputs,
		Method:  "Task.ArLink",
	}
	outfile := actions.File{
		Filename: outfile_name,
		Action:   &ar_action,
		Sem:      make(chan int, 1),
	}

	file_list[outfile_name] = &outfile
	return &outfile
}

func (t *AppCTarget) GetName() string {
	return t.Name
}

func (t *AppCTarget) GetSources() []string {
	return t.Sources
}

func (t *AppCTarget) GetResources() []string {
	return t.Resources
}

func (t *AppCTarget) GetDependencies() []Target {
	return t.Dependencies
}

func (t *AppCTarget) GetOutputFile() *actions.File {
	if file_list == nil {
		file_list = make(map[string]*actions.File)
	}

	outfile_name := strings.TrimPrefix(t.Name, "//")

	f, ok := file_list[outfile_name]
	if ok {
		return f
	}

	c_inputs := actions.MakeCObjects(t.Name, t.Sources, file_list)
	inputs := make([]*actions.File, len(c_inputs)+len(t.Dependencies))

	in_count := len(c_inputs)

	for i := range c_inputs {
		inputs[i] = c_inputs[i]
	}

	for i := range t.Dependencies {
		inputs[i+in_count] = t.Dependencies[i].GetOutputFile()
	}

	link_action := actions.Action{
		Name:    outfile_name,
		Infiles: inputs,
		Method:  "Task.LdLink",
	}

	outfile := actions.File{
		Filename: outfile_name,
		Action:   &link_action,
		Sem:      make(chan int, 1),
	}

	file_list[outfile_name] = &outfile
	return &outfile
}
