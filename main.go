package main

import (
	"os"

	"crudify/app"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
	})

	err := app.RunCliApp()
	if err != nil {
		logrus.Fatal(err)
	}
}
