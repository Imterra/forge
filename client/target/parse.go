package target

import (
	"../../log"
	"../util"
	"fmt"
	"github.com/smallfish/simpleyaml"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var target_list map[string]Target
var temp_target_list map[string]int

func MakeTarget(targetname, root, cur_dir string) Target {

	name := GetFQTN(targetname, root, cur_dir)

	if temp_target_list == nil {
		temp_target_list = make(map[string]int)
	}

	_, ok := temp_target_list[name]
	if ok {
		affected_targets := make([]string, len(temp_target_list))
		i := 0
		for k := range temp_target_list {
			affected_targets[i] = k
			i++
		}
		log.Error(
			fmt.Sprintf(
				"target %s contains circular dependencies\n(affected targets: %s)",
				name, strings.Join(affected_targets, ", ")), util.Exiter)
	}

	if target_list == nil {
		target_list = make(map[string]Target)
	}

	t, ok := target_list[name]
	if ok {
		return t
	}

	build_file := GetFile(name, root)
	temp_target_list[name] = 1
	t = ParseFile(build_file, name, root)
	delete(temp_target_list, name)
	target_list[name] = t
	return t
}

func GetFile(targetname, packageroot string) string {
	targetpath := targetname

	if IsAbs(targetname) {
		targetpath = strings.Replace(targetname, "/", packageroot, 1)
	}
	fdir := filepath.Dir(targetpath)

	return filepath.Join(fdir, "build.yaml")
}

func ParseFile(path, targetname, packageroot string) Target {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err.Error(), util.Exiter)
	}

	yaml, err := simpleyaml.NewYaml(source)
	if err != nil {
		log.Error(fmt.Sprintf("cannot parse file %s (not a valid YAML)", path), util.Exiter)
	}

	filedir := filepath.Dir(path)

	rel_targetname := filepath.Base(targetname)

	targetdata := yaml.Get(rel_targetname)

	targettype, err := targetdata.Get("type").String()
	if err != nil {
		log.Error(fmt.Sprintf("target %s does not exist", targetname), util.Exiter)
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
		log.Error("invalid type for LibCTarget", util.Exiter)
	}

	dependencies := MakeDependencies(t_data, p_root, p_cur)
	resources := GetStringArray("resources", t_data, p_root, p_cur)
	sources := GetStringArray("sources", t_data, p_root, p_cur)

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
		log.Error("invalid type for AppCTarget", util.Exiter)
	}

	dependencies := MakeDependencies(t_data, p_root, p_cur)
	resources := GetStringArray("resources", t_data, p_root, p_cur)
	sources := GetStringArray("sources", t_data, p_root, p_cur)

	return &AppCTarget{
		Name:         t_name,
		Sources:      sources,
		Resources:    resources,
		Dependencies: dependencies,
	}
}
