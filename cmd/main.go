package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/gudn/lux"
	flag "github.com/spf13/pflag"
)

var (
	projectSearchPath string
	configPath        string
	templatesPath     string
	execNginx         bool
)

func init() {
	flag.StringVarP(&projectSearchPath, "projects", "p", "projects", "path to search projects")
	flag.StringVarP(&configPath, "out", "o", "config", "path to output config files")
	flag.StringVarP(&templatesPath, "templates", "t", "templates", "path to search templates")
	flag.BoolVarP(&execNginx, "exec", "e", false, "execve to nginx master process")
	flag.Parse()
}

func main() {
	templates, err := lux.ParseTemplates(templatesPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	entries, err := os.ReadDir(projectSearchPath)
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
	l, err := lux.New(projects, templates, configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	env := os.Environ()
	root, err := filepath.Abs(l.RootConfig())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !execNginx {
		return
	}
	args := []string{"nginx", "-c", root, "-g", "daemon off;"}
	nginx, err := exec.LookPath("nginx")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	err = syscall.Exec(nginx, args, env)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
