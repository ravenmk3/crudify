package engine

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"crudify/schema/common"
	"crudify/schema/mysql"
	"crudify/utils"
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Generator struct {
	config    *ConfigModel
	tmplDir   string
	outputDir string
}

type genContext struct {
	Manifest *ManifestModel
	Vars     map[string]any
	Tables   []*common.TableSchema
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
	ctx := new(genContext)

	err := g.readManifest(ctx)
	if err != nil {
		return err
	}

	err = g.readDbSchema(ctx)
	if err != nil {
		return err
	}

	globalTpls := ctx.Manifest.GlobalTemplates
	entityTpls := ctx.Manifest.EntityTemplates

	logrus.Infof("GlobalTemplates: %d, EntityTemplates: %d, Tables: %d",
		len(globalTpls), len(entityTpls), len(ctx.Tables))

	now := time.Now()
	builtinVars := utils.Variables{
		"Db":       g.config.Database,
		"DateTime": now.Format(("2006-01-02 15:04:05")),
		"Date":     now.Format("2006-01-02"),
	}

	ctx.Vars = utils.MergeVariables(builtinVars, ctx.Manifest.Variables, g.config.Variables)

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

func (g *Generator) readManifest(ctx *genContext) error {
	manifestFile := path.Join(g.tmplDir, "manifest.yaml")
	manifest, err := ReadManifest(manifestFile)
	if err != nil {
		return err
	}
	if manifest.GlobalTemplates == nil {
		manifest.GlobalTemplates = []TemplateProps{}
	}
	if manifest.EntityTemplates == nil {
		manifest.EntityTemplates = []TemplateProps{}
	}
	ctx.Manifest = manifest
	return nil
}

func (g *Generator) readDbSchema(ctx *genContext) error {
	dbc := g.config.Database
	logrus.Infof("Reading database schema - %s:%d/%s",
		dbc.Host, dbc.Port, dbc.Database)

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

	ctx.Tables = tables
	return nil
}

func (g *Generator) renderGlobalTemplates(ctx *genContext) error {
	logrus.Info("Rendering global templates")

	for _, tplProps := range ctx.Manifest.GlobalTemplates {
		err := g.renderGlobalTemplate(ctx, &tplProps)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) renderGlobalTemplate(ctx *genContext, props *TemplateProps) error {
	logrus.Infof("Rendering global template: %s", props.File)

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

	data := &GlobalTemplateData{
		Vars:   utils.MergeVariables(ctx.Vars),
		Tables: ctx.Tables,
	}

	err = g.runGlobalScripts(ctx, props.Script, data)
	if err != nil {
		return err
	}

	outputPath, err := resolveGlobalOutputPath(props.Output, data)
	if err != nil {
		return err
	}

	return g.renderToFile(tmpl, data, outputPath)
}

func (g *Generator) runGlobalScripts(ctx *genContext, scriptFile string, data any) error {
	files := []string{}

	for _, file := range ctx.Manifest.GlobalScripts {
		files = append(files, file)
	}
	if scriptFile != "" {
		files = append(files, scriptFile)
	}

	return g.runScripts(files, "Model", data)
}

func (g *Generator) renderEntityTemplates(ctx *genContext) error {
	logrus.Info("Rendering entity templates")

	for _, tplProps := range ctx.Manifest.EntityTemplates {
		err := g.renderEntityTemplate(ctx, &tplProps)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) renderEntityTemplate(ctx *genContext, props *TemplateProps) error {
	logrus.Debugf("Rendering entity template: %s", props.File)

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

	progress, bar := NewEntityTemplateProgress(len(ctx.Tables), props.File)
	defer progress.Wait()

	for _, table := range ctx.Tables {
		err = g.renderEntityTemplateWithTable(ctx, tmpl, table, props)
		bar.Increment()
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) renderEntityTemplateWithTable(ctx *genContext, tmpl *template.Template,
	table *common.TableSchema, props *TemplateProps) error {

	logrus.Debugf("Rendering entity template: %s, %s", tmpl.Name(), table.Name)

	data := &EntityTemplateData{
		Global: &GlobalTemplateData{
			Vars:   utils.MergeVariables(ctx.Vars),
			Tables: ctx.Tables,
		},
		Vars:  utils.Variables{},
		Table: table,
	}

	err := g.runEntityScripts(ctx, props.Script, data)
	if err != nil {
		return err
	}

	outputPath, err := resolveEntityOutputPath(props.Output, data)
	if err != nil {
		return err
	}

	return g.renderToFile(tmpl, data, outputPath)
}

func (g *Generator) runEntityScripts(ctx *genContext, scriptFile string, data any) error {
	files := []string{}

	for _, file := range ctx.Manifest.EntityScripts {
		files = append(files, file)
	}
	if scriptFile != "" {
		files = append(files, scriptFile)
	}

	return g.runScripts(files, "Model", data)
}

func (g *Generator) renderToFile(tmpl *template.Template, data any, outputPath string) error {
	outputPath = path.Join(g.outputDir, outputPath)
	outputDir := path.Dir(outputPath)
	err := os.MkdirAll(outputDir, 0o755)
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		if e := f.Close(); e != nil {
			logrus.Error(e)
		}
	}(file)

	writer := bufio.NewWriter(file)
	defer func(w *bufio.Writer) {
		if e := w.Flush(); e != nil {
			logrus.Error(e)
		}
	}(writer)

	return tmpl.Execute(writer, data)
}

func (g *Generator) runScripts(scriptFiles []string, varName string, data any) error {
	if scriptFiles == nil || len(scriptFiles) <= 0 {
		return nil
	}

	scripts := []string{}

	for _, name := range scriptFiles {
		filePath := filepath.Join(g.tmplDir, name)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		scripts = append(scripts, string(content))
	}

	script := strings.Join(scripts, "\n\n")
	vm := otto.New()

	fns := &JsFunctions{}
	err := vm.Set("Utils", fns)
	if err != nil {
		return err
	}
	err = vm.Set("F", fns)
	if err != nil {
		return err
	}

	err = vm.Set(varName, data)
	if err != nil {
		return err
	}

	_, err = vm.Run(script)
	if err != nil {
		return err
	}

	return nil
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

func NewEntityTemplateProgress(total int, name string) (*mpb.Progress, *mpb.Bar) {
	progress := mpb.New(mpb.WithWidth(20))

	bar := progress.New(int64(total),
		mpb.BarStyle(),
		mpb.PrependDecorators(
			decor.CountersNoUnit("Rendering: %d/%d"),
		),
		mpb.AppendDecorators(
			decor.Name(name),
		))

	return progress, bar
}
