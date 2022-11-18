package tmpl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bit101/go-ansi"
	"github.com/bit101/ppig/conf"
	"github.com/bit101/ppig/util"
)

type Token struct {
	Name    string `json:"name"`
	Token   string `json:"token"`
	Default string `json:"default"`
	Prefix  string `json:"prefix"`
}

type TemplateDef struct {
	Name     string
	Author   string
	Language string
	Post     string
	Tokens   []Token
	SkipDir  bool
}

func GetTemplate() string {
	items := GetTemplates()
	var choiceNum int
	ok := false
	printChoices(items)
	for !ok {

		ansi.Print(ansi.Yellow, "Choice: ")
		var choice string
		_, err := fmt.Scanln(&choice)
		if err != nil {
			ansi.ClearLine()
			ansi.Println(ansi.Red, "Could not read choice")
			ansi.MoveUp(len(items) + 4)
			printChoices(items)
			ansi.ClearLine()
			continue
		}
		if strings.ToLower(choice) == "q" {
			os.Exit(0)
		}
		choice64, err := strconv.ParseInt(choice, 10, 64)
		if err != nil {
			ansi.ClearLine()
			ansi.Printf(ansi.Red, "Invalid template choice: '%s'\n", choice)
			ansi.MoveUp(len(items) + 4)
			printChoices(items)
			ansi.ClearLine()
			continue
		}
		choiceNum = int(choice64)
		if int(choiceNum) < 1 || int(choiceNum) > len(items) {
			ansi.ClearLine()
			ansi.Printf(ansi.Red, "Template choice not in range 1 to %d\n", len(items))
			ansi.MoveUp(len(items) + 4)
			printChoices(items)
			continue
		}
		ok = true
	}
	templateName := items[choiceNum-1]
	ansi.ClearLine()
	ansi.Printf(ansi.Green, "Using '%s' template\n", templateName)
	return templateName
}

func printChoices(items []string) {
	ansi.Printf(ansi.Yellow, "Choose a template (1 to %d):\n", len(items))
	for i := 0; i < len(items); i++ {
		fmt.Printf("%2d. %s\n", i+1, items[i])
	}
	fmt.Println(" Q. Quit")
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
	fmt.Printf("templateDef %+v\n", templateDef)
	return templateDef
}

func GetTokenValues(tokens []Token) map[string]string {
	tokenValues := map[string]string{}
	scanner := bufio.NewScanner(os.Stdin)
	for _, token := range tokens {
		ansi.Printf(ansi.Yellow, "%s: ", token.Name)
		fmt.Printf("%s", token.Prefix)
		if scanner.Scan() {
			tokenValues[token.Token] = token.Prefix + scanner.Text()
		}
	}
	return tokenValues
}
