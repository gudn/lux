package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gudn/lux"
	flag "github.com/spf13/pflag"
)

var (
	projectSearchPath string
	configPath        string
	templatesPath     string
)

func init() {
	flag.StringVarP(&projectSearchPath, "projects", "p", "projects", "path to search projects")
	flag.StringVarP(&configPath, "out", "o", "config", "path to output config files")
	flag.StringVarP(&templatesPath, "templates", "t", "templates", "path to search templates")
	flag.Parse()
}

func main() {
	entries, err := os.ReadDir(templatesPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	templates := make([]string, 0)
	for _, e := range entries {
		if e.Type().IsRegular() {
			name := filepath.Join(templatesPath, e.Name())
			templates = append(templates, name)
		}
	}
	entries, err = os.ReadDir(projectSearchPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	projects := make([]string, 0)
	for _, e := range entries {
		if e.Type().IsDir() {
			name := filepath.Join(projectSearchPath, e.Name())
			projects = append(projects, name)
		}
	}
	_, err = lux.New(projects, templates, configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
