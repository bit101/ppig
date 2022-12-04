// Package tmpl manages templates
package tmpl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/bit101/go-ansi"
	"github.com/bit101/ppig/conf"
	"github.com/bit101/ppig/util"
)

// Token represents a replaceable token in a template file.
type Token struct {
	Name    string `json:"name"`
	Token   string `json:"token"`
	Default string `json:"default"`
	Prefix  string `json:"prefix"`
}

// TemplateDef is the structure of a template.
type TemplateDef struct {
	Name     string
	Author   string
	Language string
	Post     string
	Tokens   []Token
	SkipDir  bool
}

// GetTemplate prompts the user go choose a template.
func GetTemplate() string {
	items := GetTemplates()
	var templateName string
	prompt := survey.Select{
		Message:  "Choose a template.",
		Options:  items,
		VimMode:  true,
		PageSize: 10,
	}
	survey.AskOne(&prompt, &templateName)
	return templateName
}

// GetTemplates gets the list of templates.
func GetTemplates() []string {
	dirItems, err := ioutil.ReadDir(conf.Config.TemplatePath)
	if err != nil {
		ansi.Println(ansi.Red, "X Unable to read templates")
		os.Exit(1)
	}

	var result []string
	for _, item := range dirItems {
		name := item.Name()
		if item.IsDir() && name != ".git" {
			result = append(result, name)
		}
	}
	return result
}

// ValidateTemplate validates that the specified template exists.
// only called if template is specified as argument to ppig.
func ValidateTemplate(template string) bool {
	templates := GetTemplates()
	for _, t := range templates {
		if t == template {
			return true
		}
	}
	return false
}

// GetTemplatePath returns the full path of the specified template.
func GetTemplatePath(templateName string) string {
	templatePath := conf.Config.TemplatePath + templateName
	if !util.DoesPathExist(templatePath) {
		ansi.Println(ansi.Red, "X Could not find template")
		os.Exit(1)
	}
	return templatePath
}

// GetTemplateConfig gets the config file for the give template.
func GetTemplateConfig(templatePath string) TemplateDef {
	var templateDef TemplateDef
	configData, err := os.ReadFile(templatePath + "/protopig.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(configData, &templateDef)
	return templateDef
}

// GetTokenValues prompts the user for a value for each token.
func GetTokenValues(tokens []Token) map[string]string {
	tokenValues := map[string]string{}
	for _, token := range tokens {
		var answer string
		prompt := survey.Input{
			Message: token.Name,
		}
		survey.AskOne(&prompt, &answer)

		tokenValues[token.Token] = token.Prefix + answer
	}
	return tokenValues
}
