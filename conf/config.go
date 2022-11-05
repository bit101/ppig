package conf

import (
	"log"
	"os"
)

type ConfigDef struct {
	TemplatePath string
}

var Config = ConfigDef{
	TemplatePath: getConfigDir() + "/protopig/protopig-templates/",
}

func getConfigDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Could not get config dir")
	}
	return dir
}
