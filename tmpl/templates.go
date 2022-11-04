package tmpl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/bit101/ppig/conf"
	"github.com/bit101/ppig/util"
)

type Token struct {
	Name     string
	Required bool
	Default  string
}

type TemplateDef struct {
	Name    string
	Author  string
	Contact string
	Type    string
	Tokens  []Token
	GoLibs  []string
}

func GetTemplate() string {
	items := GetTemplates()
	fmt.Printf("Choose a template (1 to %d):\n", len(items))
	for i := 0; i < len(items); i++ {
		fmt.Printf("%2d. %s\n", i+1, items[i])
	}
	fmt.Print("Choice: ")

	var choice int
	_, err := fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > len(items) {
		util.PrintRed("Invalid template choice")
		os.Exit(1)
	}
	return items[choice-1]
}

func GetTemplates() []string {
	dirItems, err := ioutil.ReadDir(conf.Config.TemplatePath)
	if err != nil {
		util.PrintRed("Unable to read templates")
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

func ValidateTemplate(template string) bool {
	templates := GetTemplates()
	for _, t := range templates {
		if t == template {
			return true
		}
	}
	return false
}

func GetTemplatePath(templateName string) string {
	templatePath := conf.Config.TemplatePath + templateName
	if !util.DoesPathExist(templatePath) {
		util.PrintRed("Could not find template")
		os.Exit(1)
	}
	return templatePath
}

func GetTemplateConfig(templatePath string) TemplateDef {
	var templateDef TemplateDef
	configData, err := os.ReadFile(templatePath + "/protopig.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(configData, &templateDef)
	return templateDef
}

func GetTokenValues(tokens []Token) map[string]string {
	tokenValues := map[string]string{}
	scanner := bufio.NewScanner(os.Stdin)
	for _, token := range tokens {
		fmt.Printf("%s: ", token.Name)
		if scanner.Scan() {
			tokenValues[token.Name] = scanner.Text()
		}
	}
	return tokenValues
}
