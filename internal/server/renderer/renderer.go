package renderer

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Renderer struct {
	template *template.Template
}

func New(pattern string) *Renderer {
	return &Renderer{
		template: template.Must(template.ParseGlob(pattern)),
	}
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.template.ExecuteTemplate(w, name, data)
}
