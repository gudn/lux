package lux

import (
	"os"
	"path/filepath"
	"text/template"
)

type Lux struct {
	Projects []string
	ConfigPath string
	Templates *template.Template
}

func (l *Lux) RootConfig() string {
	return filepath.Join(l.ConfigPath, "nginx.conf")
}

func (l *Lux) renderRoot() error {
	f, err := os.OpenFile(l.RootConfig(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o664)
	if err != nil {
		return err
	}
	defer f.Close()
	err = l.Templates.ExecuteTemplate(f, "nginx.conf", l)
	if err != nil {
		return err
	}
	return nil
}

func (l *Lux) renderProject(p string, result chan<- error)  {
	result <- nil
}

func (l *Lux) renderProjects() error {
	results := make(chan error, len(l.Projects))
	for _, p := range l.Projects {
		go l.renderProject(p, results)
	}
	var err error
	for range l.Projects {
		res := <- results
		if res != nil {
			err = res
		}
	}
	return err
}

func New(projects, templates []string, configPath string) (*Lux, error) {
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return nil, err
	}
	lux := &Lux{projects,configPath,tmpl}
	err = lux.renderProjects()
	if err != nil {
		return nil, err
	}
	err = lux.renderRoot()
	if err != nil {
		return nil, err
	}
	return lux, nil
}
