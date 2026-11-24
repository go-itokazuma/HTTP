package main

import (
	"fmt"
	"html/template"
)

func loadTemplates() map[string]*template.Template {
	files := map[string]string{
		"List":   "templates/list/list.html.tmpl",
		"Detail": "templates/detail/detail.html.tmpl",
		"New":    "templates/new/new.html.tmpl",
		"Edit":   "templates/edit/edit.html.tmpl",
	}

	templates := make(map[string]*template.Template)
	for name, path := range files {
		templates[name] = template.Must(template.ParseFiles(path))
	}
	return templates
}

func getTemplate(templates map[string]*template.Template, name string) (*template.Template, error) {
	tmpl, ok := templates[name]
	if !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}
	return tmpl, nil
}
