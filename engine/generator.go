package engine

import (
	"bytes"
	"log"
	"os"
	"path"
	"text/template"

	"crudify/db/common"
	"crudify/db/mysql"
	"crudify/utils"
	"github.com/sirupsen/logrus"
)

type Generator struct {
	config    *ConfigModel
	tmplDir   string
	outputDir string
}

type generationContext struct {
	Manifest   *ManifestModel
	GlobalData *GlobalTemplateData
	Tables     []common.TableSchema
}

func NewGenerator(tmplDir, outputDir, configFile string) (*Generator, error) {
	config, err := ReadConfig(configFile)
	if err != nil {
		return nil, err
	}

	g := &Generator{
		config:    config,
		tmplDir:   tmplDir,
		outputDir: outputDir,
	}
	return g, nil
}

func (g *Generator) Execute() error {
	ctx := new(generationContext)

	err := g.readManifest(ctx)
	if err != nil {
		return err
	}

	err = g.readDbSchema(ctx)
	if err != nil {
		return err
	}

	globalVars := utils.MergeVariables(ctx.Manifest.Variables, g.config.Variables)
	ctx.GlobalData = &GlobalTemplateData{
		Vars:   globalVars,
		Tables: ctx.Tables,
	}

	err = g.renderGlobalTemplates(ctx)
	if err != nil {
		return err
	}

	err = g.renderEntityTemplates(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) readManifest(ctx *generationContext) error {
	manifestFile := path.Join(g.tmplDir, "manifest.yaml")
	manifest, err := ReadManifest(manifestFile)
	if err != nil {
		return err
	}
	ctx.Manifest = manifest
	return nil
}

func (g *Generator) readDbSchema(ctx *generationContext) error {
	logrus.Info("Reading database schema")

	dbc := g.config.Database
	provider, err := mysql.NewMySqlSchemaProvider(dbc.Host, dbc.Port, dbc.Username, dbc.Password)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		e := provider.Close()
		if e != nil {
			logrus.Error(e)
		}
	}()

	tables, err := provider.GetTables(dbc.Database)
	if err != nil {
		return err
	}

	logrus.Infof("%d tables found", len(tables))
	ctx.Tables = tables
	return nil
}

func (g *Generator) renderGlobalTemplates(ctx *generationContext) error {
	for _, tplProps := range ctx.Manifest.GlobalTemplates {
		err := g.renderGlobalTemplate(ctx, &tplProps)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) renderGlobalTemplate(ctx *generationContext, props *TemplateProps) error {
	file := path.Join(g.tmplDir, props.File)
	tplBytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	content := string(tplBytes)
	tmpl, err := template.New(props.File).Parse(content)
	if err != nil {
		return err
	}

	outputPath, err := resolveGlobalOutputPath(props.Output, ctx.GlobalData)
	if err != nil {
		return err
	}

	return g.renderToFile(tmpl, ctx.GlobalData, outputPath)
}

func (g *Generator) renderEntityTemplates(ctx *generationContext) error {
	for _, tplProps := range ctx.Manifest.EntityTemplates {
		err := g.renderEntityTemplate(ctx, &tplProps)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) renderEntityTemplate(ctx *generationContext, props *TemplateProps) error {
	file := path.Join(g.tmplDir, props.File)
	tplBytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	content := string(tplBytes)
	tmpl, err := template.New(props.File).Parse(content)
	if err != nil {
		return err
	}

	for _, table := range ctx.Tables {
		err = g.renderEntityTemplateWithTable(ctx, tmpl, &table, props.Output)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) renderEntityTemplateWithTable(ctx *generationContext,
	tmpl *template.Template, table *common.TableSchema, outputPattern string) error {

	data := &EntityTemplateData{
		Global: ctx.GlobalData,
		Vars:   utils.Variables{},
		Table:  table,
	}

	outputPath, err := resolveEntityOutputPath(outputPattern, data)
	if err != nil {
		return err
	}

	return g.renderToFile(tmpl, data, outputPath)
}

func (g *Generator) renderToFile(tmpl *template.Template, data any, outputPath string) error {
	outputPath = path.Join(g.outputDir, outputPath)
	outputDir := path.Dir(outputPath)
	err := os.MkdirAll(outputDir, 0o755)
	if err != nil {
		return err
	}

	writer, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer func(w *os.File) {
		_ = w.Close()
	}(writer)

	return tmpl.Execute(writer, data)
}

func resolveGlobalOutputPath(pattern string, data *GlobalTemplateData) (string, error) {
	tmpl, err := template.New("path").Parse(pattern)
	if err != nil {
		return "", err
	}

	var arr []byte
	buf := bytes.NewBuffer(arr)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func resolveEntityOutputPath(pattern string, data *EntityTemplateData) (string, error) {
	tmpl, err := template.New("path").Parse(pattern)
	if err != nil {
		return "", err
	}

	var arr []byte
	buf := bytes.NewBuffer(arr)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
