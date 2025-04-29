package lib

import (
	"bytes"
	"html/template"
)

func RenderTemplate(path string, args any) (string, error) {
	files := []string{
		"templates/html/base.html.tmpl",
		path,
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		return "", err
	}

	w := new(bytes.Buffer)
	err = ts.ExecuteTemplate(w, "base", args)

	return w.String(), err
}
