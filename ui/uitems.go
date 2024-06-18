package ui

import "github.com/manifoldco/promptui"

type Item struct {
	Label        string
	Value        interface{}
	Details      string
	DisplayValue string //could be a function, but I am lazy
}

func (i *Item) String() string {
	return i.Label
}

const unicorn = "\U0001F984"

const (
	EXIT = "EXIT"
	BACK = "BACK"
)

// Rewrite the above as variable instantiation witt OK=&Item{Label: "OK"} pattern
var OK = &Item{Label: "OK", Details: "Confirm and proceed"}
var Back = &Item{Label: BACK, Details: "Go back to previous menu"}
var Exit = &Item{Label: EXIT, Details: "Exit the program"}

var ItemTemplate = &promptui.SelectTemplates{
	Label:    "{{ . | bold | cyan}}",
	Inactive: `{{if eq .Label "BACK"}}{{.Label | yellow}}{{else if eq .Label "EXIT"}}{{.Label | red}}{{else}}{{ .Label }}{{with .DisplayValue}}: {{.}}{{end}}{{end}}`,
	Active:   `{{if eq .Label "BACK"}}{{.Label | yellow | bold | underline}}{{else if eq .Label "EXIT"}}{{.Label | red | bold | underline}}{{else}}{{ .Label | bold | underline }}{{with .DisplayValue}}: {{. | bold}}{{end}}{{end}}`,
	Selected: "{{ .Label | red }}",
	Details:  "{{ .Details | faint }}",
}
