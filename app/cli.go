package app

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	AppName = "crudify"
)

func RunCliApp() error {
	app := NewCliApp()
	return app.Run(os.Args)
}

func NewCliApp() *cli.App {
	app := &cli.App{
		Name:        AppName,
		Usage:       AppName,
		Description: "Template based CRUD code generator",
		Commands: []*cli.Command{
			NewGenerateCommand(),
		},
	}
	return app
}

func NewGenerateCommand() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Aliases: []string{"gen"},
		Usage:   "Generate code",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug", Required: false, Value: false},
			&cli.BoolFlag{Name: "trace", Required: false, Value: false},
			&cli.StringFlag{Name: "template", Aliases: []string{"t"}, Required: false, Value: AppName + ".template"},
			&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Required: false, Value: AppName + ".output"},
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Required: false, Value: AppName + ".config.yaml"},
		},
		Action: func(ctx *cli.Context) error {
			debug := ctx.Bool("debug")
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			trace := ctx.Bool("trace")
			if trace {
				logrus.SetLevel(logrus.TraceLevel)
			}
			tmplDir := ctx.String("template")
			outputDir := ctx.String("output")
			configFile := ctx.String("config")
			return ExecGenerate(tmplDir, outputDir, configFile)
		},
	}
}
