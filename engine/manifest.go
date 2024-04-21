package engine

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TemplateProps struct {
	File   string `yaml:"file"`
	Script string `yaml:"script"`
	Output string `yaml:"output"`
}

type ManifestModel struct {
	Variables       map[string]any  `yaml:"variables"`
	GlobalScripts   []string        `yaml:"global-scripts"`
	GlobalTemplates []TemplateProps `yaml:"global-templates"`
	EntityScripts   []string        `yaml:"entity-scripts"`
	EntityTemplates []TemplateProps `yaml:"entity-templates"`
}

func ReadManifest(filename string) (*ManifestModel, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := new(ManifestModel)
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
