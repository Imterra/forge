package target

import (
	"github.com/smallfish/simpleyaml"
	"path/filepath"
	"strings"
)

func GetStringArray(key string, data *simpleyaml.Yaml,
	packageroot, curwd string) []string {
	value := data.Get(key)

	if value == nil {
		return make([]string, 0)
	}

	value_arr, _ := value.Array()

	string_array := make([]string, len(value_arr))

	for i := 0; i < len(value_arr); i++ {
		str_val, _ := value.GetIndex(i).String()
		if !filepath.IsAbs(str_val) {
			str_val = GetFQTN(str_val, packageroot, curwd)
			//str_val = filepath.Join(curwd, str_val)
		}
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
