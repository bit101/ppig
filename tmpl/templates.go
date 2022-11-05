package tmpl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/bit101/go-ansi"
	"github.com/bit101/ppig/conf"
	"github.com/bit101/ppig/util"
)

type Token struct {
	Name    string
	Token   string
	Default string
	Prefix  string
}

type TemplateDef struct {
	Name     string
	Author   string
	Contact  string
	Language string
	Post     string
	Tokens   []Token
	GoLibs   []string
}

func GetTemplate() string {
	items := GetTemplates()
	ansi.Printf(ansi.Brown, "Choose a template (1 to %d):\n", len(items))
	for i := 0; i < len(items); i++ {
		fmt.Printf("%2d. %s\n", i+1, items[i])
	}
	ansi.Print(ansi.Brown, "Choice: ")

	var choice string
	_, err := fmt.Scanln(&choice)
	if err != nil {
		ansi.Println(ansi.Red, "Could not read choice")
		os.Exit(1)
	}
	choice64, err := strconv.ParseInt(choice, 10, 64)
	if err != nil {
		ansi.Printf(ansi.Red, "Invalid template choice: '%s'\n", choice)
		os.Exit(1)
	}
	choiceNum := int(choice64)
	if int(choiceNum) < 1 || int(choiceNum) > len(items) {
		ansi.Printf(ansi.Red, "Template choice not in range 1 to %d\n", len(items))
		os.Exit(1)
	}
	templateName := items[choiceNum-1]
	ansi.Printf(ansi.Green, "Using '%s' template\n", templateName)
	return templateName
}

func GetTemplates() []string {
	dirItems, err := ioutil.ReadDir(conf.Config.TemplatePath)
	if err != nil {
		ansi.Println(ansi.Red, "Unable to read templates")
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
		ansi.Println(ansi.Red, "Could not find template")
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
		ansi.Printf(ansi.Brown, "%s: ", token.Name)
		fmt.Printf("%s", token.Prefix)
		if scanner.Scan() {
			tokenValues[token.Token] = token.Prefix + scanner.Text()
		}
	}
	return tokenValues
}
