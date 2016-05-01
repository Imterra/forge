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
	GetOutputFile() *actions.GeneratedFile
}

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

func (t *LibCTarget) GetOutputFile() *actions.GeneratedFile {
	inputs := actions.MakeCObjects(t.Name, t.Sources)
	ar_action := actions.Action{
		Name:    strings.TrimPrefix(t.Name, "//"),
		Infiles: inputs,
		Method:  "Task.ArLink",
	}
	outfile := actions.GeneratedFile{
		Filename: strings.TrimPrefix(t.Name, "//") + ".a",
		Action:   &ar_action,
	}
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

func (t *AppCTarget) GetOutputFile() *actions.GeneratedFile {
	c_inputs := actions.MakeCObjects(t.Name, t.Sources)
	inputs := make([]actions.File, len(c_inputs)+len(t.Dependencies))

	in_count := len(c_inputs)

	for i := range c_inputs {
		inputs[i] = c_inputs[i]
	}

	for i := range t.Dependencies {
		inputs[i+in_count] = t.Dependencies[i].GetOutputFile()
	}

	link_action := actions.Action{
		Name:    strings.TrimPrefix(t.Name, "//"),
		Infiles: inputs,
		Method:  "Task.LdLink",
	}

	outfile := actions.GeneratedFile{
		Filename: strings.TrimPrefix(t.Name, "//"),
		Action:   &link_action,
	}
	return &outfile
}
