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

func New(path string, funcMap *template.FuncMap) *Renderer {
	baseTemplate := template.New("").Funcs(*funcMap)

	tmpls := map[string]*template.Template{
		constants.ChannelsTemplate: template.Must(template.Must(baseTemplate.Clone()).ParseFiles(
			filepath.Join(path, constants.BaseTemplate),
			filepath.Join(path, constants.ChannelsTemplate),
		)),
		constants.ItemsTemplate: template.Must(template.Must(baseTemplate.Clone()).ParseFiles(
			filepath.Join(path, constants.BaseTemplate),
			filepath.Join(path, constants.ItemsTemplate),
		)),
		constants.ItemTemplate: template.Must(template.Must(baseTemplate.Clone()).ParseFiles(
			filepath.Join(path, constants.BaseTemplate),
			filepath.Join(path, constants.ItemTemplate),
		)),
		constants.MessageTemplate: template.Must(template.Must(baseTemplate.Clone()).ParseFiles(
			filepath.Join(path, constants.BaseTemplate),
			filepath.Join(path, constants.MessageTemplate),
		)),
	}

	return &Renderer{templates: tmpls}
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("template %v not found", name)
	}

	return tmpl.ExecuteTemplate(w, "base", data)
}
