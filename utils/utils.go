package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

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
		PrintRed("Invalid template choice")
		os.Exit(1)
	}
	return items[choice-1]
}

func GetTemplates() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dirName := home + "/Templates/tinpig-templates/"
	dirItems, err := ioutil.ReadDir(dirName)
	if err != nil {
		PrintRed("Unable to read templates")
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

func GetProjectPath() string {
	fmt.Print("Create project in directory: ")
	var projectPath string
	fmt.Scanln(&projectPath)
	// todo validate project identifier
	if projectPath == "" {
		PrintRed("No directory specified")
		os.Exit(1)
	}
	if DoesPathExist(projectPath) {
		PrintRed(fmt.Sprintf("Something already here with the name '%s'", projectPath))
		os.Exit(1)
	}
	return projectPath
}

func DoesPathExist(projectPath string) bool {
	_, err := os.Stat(projectPath)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func PrintRed(msg string) {
	fmt.Printf("\033[1;31m%s\033[0m\n", msg)
}
