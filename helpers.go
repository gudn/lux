package lux

import (
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/net/idna"
)

func join(args ...string) string {
	return filepath.Join(args...)
}

func punycode(arg string) (string, error) {
	return idna.ToASCII(arg)
}

func ParseTemplates(templatesPath string) (*template.Template, error) {
	fmap := template.FuncMap{"join": join, "punycode": punycode}
	tmpl := template.New("").Funcs(fmap)
	entries, err := os.ReadDir(templatesPath)
	if err != nil {
		return nil, err
	}
	tmplfiles := make([]string, 0)
	for _, e := range entries {
		if e.Type().IsRegular() {
			tmplfiles = append(tmplfiles, filepath.Join(templatesPath, e.Name()))
		}
	}
	return tmpl.ParseFiles(tmplfiles...)
}
