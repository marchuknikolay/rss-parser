package renderer

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/marchuknikolay/rss-parser/internal/server/templates/constants"
)

type Renderer struct {
	templates map[string]*template.Template
}

func New(path string, funcMap *template.FuncMap) (*Renderer, error) {
	baseTemplate := template.New("").Funcs(*funcMap)

	loadTemplate := func(filenames ...string) (*template.Template, error) {
		cloned, err := baseTemplate.Clone()
		if err != nil {
			return nil, fmt.Errorf("failed to clone a base template: %w", err)
		}

		return cloned.ParseFiles(filenames...)
	}

	tmpls := make(map[string]*template.Template)

	var err error

	if tmpls[constants.ChannelsTemplate], err = loadTemplate(
		filepath.Join(path, constants.BaseTemplate),
		filepath.Join(path, constants.ChannelsTemplate)); err != nil {
		return nil, fmt.Errorf("load channels template: %w", err)
	}

	if tmpls[constants.ItemsTemplate], err = loadTemplate(
		filepath.Join(path, constants.BaseTemplate),
		filepath.Join(path, constants.ItemsTemplate)); err != nil {
		return nil, fmt.Errorf("load items template: %w", err)
	}

	if tmpls[constants.ItemTemplate], err = loadTemplate(
		filepath.Join(path, constants.BaseTemplate),
		filepath.Join(path, constants.ItemTemplate)); err != nil {
		return nil, fmt.Errorf("load item template: %w", err)
	}

	if tmpls[constants.MessageTemplate], err = loadTemplate(
		filepath.Join(path, constants.BaseTemplate),
		filepath.Join(path, constants.MessageTemplate)); err != nil {
		return nil, fmt.Errorf("load message template: %w", err)
	}

	return &Renderer{templates: tmpls}, nil
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("template %v not found", name)
	}

	return tmpl.ExecuteTemplate(w, "base", data)
}
