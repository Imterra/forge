package target

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/smallfish/simpleyaml"
	"io/ioutil"
)

func ParseFile(filepath, targetname string) Target {
	source, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	yaml, err := simpleyaml.NewYaml(source)
	if err != nil {
		panic(err)
	}

	spew.Dump(yaml)
	return nil
}
