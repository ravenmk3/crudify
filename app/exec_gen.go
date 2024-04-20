package app

import (
	"crudify/engine"
	"github.com/sirupsen/logrus"
)

func ExecGenerate(tmplDir, outputDir, configFile string) error {
	logrus.Info("Generation started")
	logrus.Infof("Template directory: %s", tmplDir)
	logrus.Infof("Output directory: %s", outputDir)
	logrus.Infof("Config file: %s", configFile)

	generator, err := engine.NewGenerator(tmplDir, outputDir, configFile)
	if err != nil {
		return err
	}
	return generator.Execute()
}
