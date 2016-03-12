package actions

import (
	"fmt"
	"strings"
)

type LibCAction struct {
	Name    string
	Infiles []File
}

func (a *LibCAction) GetCmd() string {
	filenames_list := make([]string, len(a.Infiles))
	for i := 0; i < len(filenames_list); i++ {
		filenames_list[i] = a.Infiles[i].GetFullPath()
	}

	return fmt.Sprintf(
		"gcc -W -Wall -Wextra -o %s -c %s",
		a.GetOutFileName(), strings.Join(filenames_list, " "))
}

func (a *LibCAction) GetOutFileName() string {
	return fmt.Sprintf("%s.o", a.Name)
}

func (a *LibCAction) GetInfiles() []File {
	return a.Infiles
}

func (a *LibCAction) GetName() string {
	return a.Name
}

type AppCAction struct {
	Name    string
	Infiles []File
}

func (a *AppCAction) GetCmd() string {
	filenames_list := make([]string, len(a.Infiles))
	for i := 0; i < len(filenames_list); i++ {
		filenames_list[i] = a.Infiles[i].GetFullPath()
	}

	return fmt.Sprintf(
		"gcc -W -Wall -Wextra -o %s %s",
		a.GetOutFileName(), strings.Join(filenames_list, " "))
}

func (a *AppCAction) GetOutFileName() string {
	return a.Name
}

func (a *AppCAction) GetInfiles() []File {
	return a.Infiles
}

func (a *AppCAction) GetName() string {
	return a.Name
}
