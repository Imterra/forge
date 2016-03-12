package target

import (
	"../actions"
)

type Target interface {
	GetName() string
	GetSources() []string
	GetResources() []string
	GetDependencies() []Target
	GetAction() actions.Action
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

func (t *LibCTarget) GetAction() actions.Action {
	infiles := GetInFiles(t)
	action := actions.LibCAction{Name: t.GetName(), Infiles: infiles}

	return &action
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

func (t *AppCTarget) GetAction() actions.Action {
	infiles := GetInFiles(t)
	action := actions.AppCAction{Name: t.GetName(), Infiles: infiles}

	return &action
}
