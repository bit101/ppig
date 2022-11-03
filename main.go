package main

import (
	"fmt"

	flag "github.com/spf13/pflag"

	"github.com/bit101/ppig/project"
	"github.com/bit101/ppig/utils"
)

func main() {
	templateFlag := flag.StringP("template", "t", "", "The template to use for the project.")
	pathFlag := flag.StringP("path", "p", "", "The path to create the project in. Must not already exist")
	configFlag := flag.BoolP("config", "c", false, "Configure protopig options.")
	helpFlag := flag.BoolP("help", "h", false, "Get help")
	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()

	if *helpFlag {
		fmt.Println("Protopig usage:\n\nppig [flags]\n\nflags:")
		flag.PrintDefaults()
		return
	}

	if *configFlag {
		fmt.Println("configuring....")
		return
	}

	templateName := *templateFlag
	projectPath := *pathFlag

	if templateName == "" {
		templateName = utils.GetTemplate()
	}

	if projectPath == "" {
		projectPath = utils.GetProjectPath()
	}

	project.CreateProject(projectPath, templateName)

}
