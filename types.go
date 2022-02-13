package lux

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Project map[string]interface{}

func (p Project) TemplateName() string {
	value, ok := p["template"]
	if ok {
		if s, ok := value.(string); ok {
			return s
		}
	}
	return "default"
}

func loadYaml(p string) (Project, error) {
	result := make(map[string]interface{})
	f, err := os.ReadFile(filepath.Join(p, "lux.yaml"))
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(f, result)
	return result, err
}
