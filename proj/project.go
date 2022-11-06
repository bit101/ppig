package proj

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bit101/go-ansi"
	"github.com/bit101/ppig/tmpl"
	"github.com/bit101/ppig/util"
	"github.com/otiai10/copy"
)

type Project struct {
	TemplateName string
	ProjectPath  string
	TemplatePath string
	Template     tmpl.TemplateDef
	TokenValues  map[string]string
}

func NewProject(projectPath, templateName string) *Project {
	p := &Project{}
	p.getTemplate(templateName)
	p.getProjectPath(projectPath)
	p.TemplatePath = tmpl.GetTemplatePath(p.TemplateName)
	p.Template = tmpl.GetTemplateConfig(p.TemplatePath)
	p.TokenValues = tmpl.GetTokenValues(p.Template.Tokens)
	return p
}

func (p *Project) Build() {
	copy.Copy(p.TemplatePath, p.ProjectPath)
	p.replaceTokens()
	p.finalize()
	ansi.Printf(ansi.Green, "A '%s' project has been created in the '%s' directory!", p.TemplateName, p.ProjectPath)
}

func (p *Project) getProjectPath(projectPath string) {
	if projectPath != "" {
		p.ProjectPath = projectPath
		return
	}
	ansi.Print(ansi.Yellow, "Create project in directory: ")
	fmt.Scanln(&projectPath)
	if projectPath == "" {
		ansi.Println(ansi.Red, "No directory specified")
		os.Exit(1)
	}
	if util.DoesPathExist(projectPath) {
		ansi.Printf(ansi.Red, "Something already here with the name '%s'\n", projectPath)
		os.Exit(1)
	}
	p.ProjectPath = projectPath
}

func (p *Project) getTemplate(templateName string) {
	if templateName == "" {
		p.TemplateName = tmpl.GetTemplate()
	} else if !tmpl.ValidateTemplate(templateName) {
		ansi.Printf(ansi.Red, "'%s' is not a valid template name.\n", templateName)
		p.TemplateName = tmpl.GetTemplate()
	} else {
		p.TemplateName = templateName
	}
}

func (p *Project) replaceTokens() {
	filepath.Walk(p.ProjectPath, func(path string, info os.FileInfo, err error) error {
		if info.Name() == "protopig.json" {
			os.Remove(path)
			// todo error check
			return nil
		}

		if info.IsDir() {
			return nil
		}

		var outContent string
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		didReplace := false
		for key, value := range p.TokenValues {
			outContent = strings.ReplaceAll(string(content), "${"+key+"}", value)
			didReplace = true
		}
		if didReplace && string(content) != outContent {
			os.WriteFile(path, []byte(outContent), 0644)
		}
		return nil
	})
}

func (p *Project) finalize() {
	if p.Template.Post != "" {
		cmd := exec.Command("./" + p.Template.Post)
		cmd.Dir = p.ProjectPath
		err := cmd.Run()
		if err != nil {
			log.Fatal("error: ", err)
		}
	}
}
