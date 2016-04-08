package actions

import (
	"../objectstorage"
	"fmt"
	"strings"
)

type LibCAction struct {
	Name    string
	Infiles []File
	Storage objectstorage.Storage
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
	outfilename := fmt.Sprintf("%s.o", a.Name)
	return a.Storage.GetFilePath(outfilename, GetSourceChecksum(a.Infiles))
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
	Storage objectstorage.Storage
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
	outfilename := a.Name
	return a.Storage.GetFilePath(outfilename, GetSourceChecksum(a.Infiles))
}

func (a *AppCAction) GetInfiles() []File {
	return a.Infiles
}

func (a *AppCAction) GetName() string {
	return a.Name
}
