package main

import (
	"fmt"

	"github.com/cyp57/uploadapi/config"
	"github.com/cyp57/uploadapi/pkg/mongodb"
	"github.com/cyp57/uploadapi/route"
	"github.com/cyp57/uploadapi/setting"
)

const (
	PathEnv  = "config/.env"
	PathIni  = "config/app.ini"
	PathYaml = "config"
)

func main() {
	config := config.LoadConfig(PathEnv, PathYaml)
	fmt.Println("config = =", config)
	fmt.Println("config = =", config.App())
	fmt.Println("config = =", config.Db())
	setting.InitIni(PathIni)
	mongodb.MongoDbConnect(config.Db())
	route.InitRoute(config.App())
}
