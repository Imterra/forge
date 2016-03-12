package util

import (
	"../actions"
	"fmt"
	"path/filepath"
)

func NormalizePath(path string) (string, string) {
	clean_path := filepath.Clean(path)
	full_path, _ := filepath.Abs(clean_path)
	return filepath.Base(full_path), full_path
}

func PrintAllActions(action actions.Action) {
	infiles := action.GetInfiles()
	for i := 0; i < len(infiles); i++ {
		origin := infiles[i].GetOrigin()
		if origin != nil {
			PrintAllActions(origin)
		}
	}

	fmt.Println(action.GetCmd())
}
