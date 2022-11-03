package project

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bit101/ppig/utils"
)

type Token struct {
	Name     string
	Required bool
	Default  string
}

type TemplateConfig struct {
	Name        string
	Author      string
	Contact     string
	Tokens      []Token
	Ignore      []string
	PostMessage string
}

func CreateProject(projectPath, templateName string) {
	templatePath := getTemplatePath(templateName)
	templateConfig := getTemplateConfig(templatePath)
	tokenValues := getTokenValues(templateConfig.Tokens)
	copyFiles(templatePath, tokenValues)
}

func getTemplatePath(templateName string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	templatePath := home + "/Templates/tinpig-templates/" + templateName
	if !utils.DoesPathExist(templatePath) {
		utils.PrintRed("Could not find template")
		os.Exit(1)
	}
	return templatePath
}

func getTemplateConfig(templatePath string) TemplateConfig {
	var templateConfig TemplateConfig
	configData, err := os.ReadFile(templatePath + "/tinpig.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(configData, &templateConfig)
	return templateConfig
}

func getTokenValues(tokens []Token) map[string]string {
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

func copyFiles(templatePath string, tokenValues map[string]string) {
	templateFiles, err := ioutil.ReadDir(templatePath)
	if err != nil {
		utils.PrintRed("Unable to read template directory")
		os.Exit(1)
	}

	for _, file := range templateFiles {
		if file.IsDir() {
			copyFiles(templatePath+"/"+file.Name(), tokenValues)
		} else if file.Name() != "tinpig.json" {
			contents, err := os.ReadFile(templatePath + "/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			for key, value := range tokenValues {
				contentStr := strings.ReplaceAll(string(contents), "${"+key+"}", value)
				// todo: make the directories and write the files
				fmt.Println(file.Name())
				fmt.Println(contentStr)
			}
		}
	}
}
