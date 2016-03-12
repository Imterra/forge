package target

import (
	"../actions"
	"../util"
)

func GetInFiles(t Target) []actions.File {
	sources := t.GetSources()
	resources := t.GetResources()
	dependencies := t.GetDependencies()

	infiles := make(
		[]actions.File, len(sources)+len(resources)+len(dependencies))

	max := len(sources)
	for i := 0; i < max; i++ {
		filename, fullpath := util.NormalizePath(sources[i])
		infiles[i] = &actions.SourceFile{Name: filename, FullPath: fullpath}
	}

	for i := 0; i < len(resources); i++ {
		filename, fullpath := util.NormalizePath(resources[i])
		infiles[i+max] = &actions.SourceFile{Name: filename, FullPath: fullpath}
	}
	max = max + len(resources)

	for i := 0; i < len(dependencies); i++ {
		a := dependencies[i].GetAction()
		filename, fullpath := util.NormalizePath(a.GetOutFileName())
		f := actions.GeneratedFile{Name: filename, FullPath: fullpath, Origin: a}
		infiles[i+max] = &f
	}

	return infiles
}
