package lux

import (
	"os"
	"path/filepath"
	"text/template"
)

type Lux struct {
	Projects   []string
	ConfigPath string
	Templates  *template.Template
}

func (l *Lux) RootConfig() string {
	return filepath.Join(l.ConfigPath, "nginx.conf")
}

func (l *Lux) renderRoot(projects []string) error {
	configPath, err := filepath.Abs(l.ConfigPath)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(l.RootConfig(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o664)
	if err != nil {
		return err
	}
	defer f.Close()
	data := map[string]interface{}{"ConfigPath": configPath, "ProjectsConfigs": projects}
	err = l.Templates.ExecuteTemplate(f, "nginx.conf", data)
	if err != nil {
		return err
	}
	return nil
}

func (l *Lux) renderProject(p string, errors chan<- error, results chan<- string) {
	data, err := loadYaml(p)
	if err != nil {
		errors <- err
		return
	}
	data["root"], err = filepath.Abs(p)
	if err != nil {
		errors <- err
		return
	}
	name := filepath.Base(p)
	path := filepath.Join(l.ConfigPath, "projects", name)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o664)
	if err != nil {
		errors <- err
		return
	}
	defer f.Close()
	err = l.Templates.ExecuteTemplate(f, data.TemplateName(), data)
	if err != nil {
		errors <- err
		return
	}
	results <- f.Name()
}

func (l *Lux) renderProjects() ([]string, error) {
	err := os.MkdirAll(filepath.Join(l.ConfigPath, "projects"), 0o776)
	if err != nil {
		return nil, err
	}
	errors := make(chan error, len(l.Projects))
	results := make(chan string, len(l.Projects))
	for _, p := range l.Projects {
		go l.renderProject(p, errors, results)
	}
	err = nil
	configs := make([]string, 0, len(l.Projects))
	for range l.Projects {
		select {
		case e := <-errors:
			err = e
		case p := <-results:
			configs = append(configs, p)
		}
	}
	return configs, err
}

func New(projects []string, templates *template.Template, configPath string) (*Lux, error) {
	lux := &Lux{projects, configPath, templates}
	configs, err := lux.renderProjects()
	if err != nil {
		return nil, err
	}
	err = lux.renderRoot(configs)
	if err != nil {
		return nil, err
	}
	return lux, nil
}
