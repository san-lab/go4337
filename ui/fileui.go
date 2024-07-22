package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

// Define an interface Item with Name and Size
type FileItem interface {
	Name() string
	IsDir() bool
}

// Define ExtraItem with Name and Size
type ExtraItem struct {
	name string
	size int64
}

func (ei *ExtraItem) Name() string {
	return ei.name
}

func (ei *ExtraItem) IsDir() bool {
	return false
}

func (ei ExtraItem) String() string {
	return ei.name
}

const selectFileLabel = "Select a file"

func SelectFileFromFS(suffix string) (filebytes []byte, err error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ .| bold}}?",
		Active:   "{{ if .IsDir }}{{.Name | cyan | bold | underline}}{{else}}{{.Name | bold | underline}}{{end}} ",
		Inactive: `{{ if .IsDir }}{{.Name | cyan }}{{else if eq .Name "BACK"}}{{.Name | yellow}}{{else}}{{.Name}}{{end}}`,
		Selected: "{{.Name | red }}",
	}

	//Create a template that makes the text bold

	for {
		prompt := promptui.Select{
			Label:     selectFileLabel,
			Templates: templates,
		}
		dir, err := os.ReadDir(".")
		if err != nil {
			return nil, err
		}
		fi := []FileItem{}
		for _, file := range dir {
			if strings.HasSuffix(file.Name(), suffix) {
				fi = append([]FileItem{file}, fi...)
			}
			if file.IsDir() {
				fi = append(fi, file)
			}
		}
		fi = append(fi, &ExtraItem{name: BACK})
		fi = append([]FileItem{&ExtraItem{name: ".."}}, fi...)

		prompt.Items = fi
		prompt.Size = len(fi)
		if prompt.Size > 10 {
			prompt.Size = 10
		}
		i, _, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		if i == 0 || fi[i].IsDir() {
			os.Chdir(fi[i].Name())
			continue
		}
		if fi[i].Name() == BACK {
			return nil, fmt.Errorf("Canceled")
		}

		//Try to read the file
		return os.ReadFile(fi[i].Name())

	}
}
