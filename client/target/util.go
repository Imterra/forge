package target

import (
	"../actions"
	"../util"
	"github.com/smallfish/simpleyaml"
	"strings"
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

func GetStringArray(key string, data *simpleyaml.Yaml) []string {
	value := data.Get(key)

	if value == nil {
		return make([]string, 0)
	}

	value_arr, _ := value.Array()

	string_array := make([]string, len(value_arr))

	for i := 0; i < len(value_arr); i++ {
		str_val, _ := value.GetIndex(i).String()
		string_array[i] = str_val
	}

	return string_array
}

func GetFQTN(target, packageroot, curdir string) string {
	if IsAbs(target) {
		return target
	}

	var curpkg string

	if strings.HasPrefix(curdir, packageroot) {
		curpkg = strings.Replace(curdir, packageroot, "/", 1)
	} else {
		curpkg = curdir
	}

	arr := make([]string, 2)
	arr[0] = curpkg
	arr[1] = target

	return strings.Join(arr, "/")
}

func IsAbs(target string) bool {
	return strings.HasPrefix(target, "//")
}
