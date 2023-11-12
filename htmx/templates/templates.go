package templates

import (
	"fmt"
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Template struct {
	RootDir string
}

func (t *Template) getExtraTemplates(name string) []string {
	var tmps = map[string][]string{
		"about.html": {t.RootDir + "templates/about.html"},
		"index.html": {t.RootDir + "templates/index.html"},
		"user.html":  {t.RootDir + "templates/user.html"},
	}
	return tmps[name]
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	tmp := template.Must(template.ParseGlob(t.RootDir + "templates/shared/*.html"))

	f := t.getExtraTemplates(name)
	if len(f) > 0 {
		if err := template.Must(tmp.ParseFiles(f...)).ExecuteTemplate(w, "base", data); err != nil {
			fmt.Println("Error", err)
			return err
		}
	} else {
		if err := tmp.ExecuteTemplate(w, name, data); err != nil {
			fmt.Println("Error", err)
			return err
		}
	}

	return nil
}
