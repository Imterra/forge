package target

import (
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	"github.com/smallfish/simpleyaml"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var target_list map[string]Target

func MakeTarget(targetname, root, cur_dir string) Target {

	name := GetFQTN(targetname, root, cur_dir)

	fmt.Printf("[DBG] TN: %v, FQTN: %v\n", targetname, name)

	if target_list == nil {
		target_list = make(map[string]Target)
	}

	t, ok := target_list[name]
	if ok {
		return t
	}

	build_file := GetFile(name, root)
	t = ParseFile(build_file, name, root)
	target_list[name] = t
	return t
}

func GetFile(targetname, packageroot string) string {
	targetpath := targetname

	if IsAbs(targetname) {
		targetpath = strings.Replace(targetname, "/", packageroot, 1)
	}
	fdir := filepath.Dir(targetpath)

	fmt.Printf("[DBG] TN: %v, BF: %v\n", targetname, filepath.Join(fdir, "build.yaml"))

	return filepath.Join(fdir, "build.yaml")
}

func ParseFile(path, targetname, packageroot string) Target {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	yaml, err := simpleyaml.NewYaml(source)
	if err != nil {
		panic(err)
	}

	filedir := filepath.Dir(path)

	rel_targetname := filepath.Base(targetname)

	targetdata := yaml.Get(rel_targetname)

	targettype, err := targetdata.Get("type").String()
	if err != nil {
		panic(err)
	}

	var target Target
	switch targettype {
	case "lib_c":
		target = MakeLibCTarget(targetname, targetdata, packageroot, filedir)
	case "app_c":
		target = MakeAppCTarget(targetname, targetdata, packageroot, filedir)
	}

	return target
}

func MakeDependencies(t_data *simpleyaml.Yaml, p_root, p_cur string) []Target {
	t_deps := t_data.Get("dependencies")

	var t_deps_count int

	if t_deps != nil {
		t_deps_arr, _ := t_deps.Array()
		t_deps_count = len(t_deps_arr)
	} else {
		t_deps_count = 0
	}

	dependencies := make([]Target, t_deps_count)

	for i := 0; i < t_deps_count; i++ {
		t_dep_str, _ := t_deps.GetIndex(i).String()
		dependencies[i] = MakeTarget(t_dep_str, p_root, p_cur)
	}

	return dependencies
}

func MakeLibCTarget(t_name string, t_data *simpleyaml.Yaml, p_root, p_cur string) *LibCTarget {
	t_type, _ := t_data.Get("type").String()

	if t_type != "lib_c" {
		panic("Invalid type for LibCTarget!")
	}

	fmt.Printf("\n[DBG] MakeLibC: p_root: %v, p_cur: %v\n\n", p_root, p_cur)

	dependencies := MakeDependencies(t_data, p_root, p_cur)
	resources := GetStringArray("resources", t_data, p_cur)
	sources := GetStringArray("sources", t_data, p_cur)

	return &LibCTarget{
		Name:         t_name,
		Sources:      sources,
		Resources:    resources,
		Dependencies: dependencies,
	}
}

func MakeAppCTarget(t_name string, t_data *simpleyaml.Yaml, p_root, p_cur string) *AppCTarget {
	t_type, _ := t_data.Get("type").String()

	if t_type != "app_c" {
		panic("Invalid type for AppCTarget!")
	}

	fmt.Printf("[DBG] MakeAppC: p_root: %v, p_cur: %v\n", p_root, p_cur)

	dependencies := MakeDependencies(t_data, p_root, p_cur)
	resources := GetStringArray("resources", t_data, p_cur)
	sources := GetStringArray("sources", t_data, p_cur)

	return &AppCTarget{
		Name:         t_name,
		Sources:      sources,
		Resources:    resources,
		Dependencies: dependencies,
	}
}
