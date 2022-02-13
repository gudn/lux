package lux

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
