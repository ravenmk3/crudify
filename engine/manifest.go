package engine

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TemplateProps struct {
	File   string `yaml:"file"`
	Output string `yaml:"output"`
}

type ManifestModel struct {
	Variables       map[string]any  `yaml:"variables"`
	GlobalTemplates []TemplateProps `yaml:"global-templates"`
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
