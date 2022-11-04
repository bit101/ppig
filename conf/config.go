package conf

import (
	"log"
	"os"
)

type ConfigDef struct {
	TemplatePath string
	GoModBase    string
}

var Config = ConfigDef{
	TemplatePath: getConfigDir() + "/protopig/protopig-templates/",
	GoModBase:    "github.com/bit101/",
}

func getConfigDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Could not get config dir")
	}
	return dir
}
