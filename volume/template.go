package volume

import (
    "html/template"
)

/*

Template

 */

type Template struct {
    *template.Template
    text string
}

func( tmpl *Template ) Text() string {
    return tmpl.text
}