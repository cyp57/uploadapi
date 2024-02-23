package setting

import (
	"log"

	"github.com/go-ini/ini"
)

type ApiGroup struct {
}

var ApiGroupSetting = &ApiGroup{}

type ApiEndpoint struct {
	LogEndpoint string
}

var ApiEndpointSetting = &ApiEndpoint{}

type Collection struct {
	GridFsCollection string
	FileserverCollection string
}

var CollectionSetting = &Collection{}

var cfg *ini.File

// initialize .ini
func InitIni(iniPath string) {

	var err error
	cfg, err = ini.Load(iniPath)
	if err != nil {
		log.Fatalln(err)
	}

	mapTo("apiGroup", ApiGroupSetting)
	mapTo("apiEndpoint", ApiEndpointSetting)
	mapTo("collection", CollectionSetting)

}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalln(err)
	}
}
