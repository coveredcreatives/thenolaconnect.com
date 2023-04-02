package template

import (
	"os"
	"path/filepath"
)

func OrderSheetTemplate() (s string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	os_template, err := os.ReadFile(filepath.Join(wd, "template", "ordersheet.html"))
	if err != nil {
		return
	}
	s = string(os_template)
	return
}
