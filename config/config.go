package config

import (
	"log"

	"github.com/cyp57/uploadapi/utils"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type app struct {
	httpPort         string
	httpRequestLimit int
	serviceName      string
	basePath         string
	rootFile         string
}

type db struct {
	dbUser     string
	dbPassword string
	dbHost     string
	dbName     string
}

type config struct {
	app *app
	db  *db
}

type IAppConfig interface {
	HTTPPort() string
	ServiceName() string
	HttpRequestLimit() int
	BasePath() string
	RootFile() string
}

type IDbConfig interface {
	DbUser() string
	DbPassword() string
	DbHost() string
	DbName() string
}

func (c *config) App() IAppConfig {
	return c.app
}

func (a *app) HTTPPort() string {
	return a.httpPort
}
func (a *app) ServiceName() string {
	return a.serviceName
}
func (a *app) HttpRequestLimit() int {
	return a.httpRequestLimit
}
func (a *app) BasePath() string {
	return a.basePath
}
func (a *app) RootFile() string {
	return a.rootFile
}

func (c *config) Db() IDbConfig {
	return c.db
}

func (d *db) DbUser() string {
	return d.dbUser
}

func (d *db) DbPassword() string {
	return d.dbPassword
}
func (d *db) DbHost() string {
	return d.dbHost
}
func (d *db) DbName() string {
	return d.dbName
}

type IConfig interface {
	App() IAppConfig
	Db() IDbConfig
}

func LoadConfig(envPath, yamlPath string) IConfig {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalln(err)
	}
	mode := utils.GetEnv("mode")

	var viperYaml = viper.New()
	viperYaml.SetConfigName(mode)

	viperYaml.SetConfigType("yaml")
	viperYaml.AddConfigPath(yamlPath)

	err = viperYaml.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	} else {
		utils.SetViperYaml(viperYaml)
	}

	return &config{&app{httpPort: utils.GetYaml("HTTPPort"),
		serviceName:      utils.GetYaml("ServiceName"),
		basePath:         utils.GetYaml("BasePath"),
		rootFile: utils.GetYaml("RootFile"),
		httpRequestLimit: utils.GetYamlInt("HttpRequestLimit")},
		
		&db{dbUser: utils.GetYaml("DBUser"),
			dbPassword: utils.GetYaml("DBPassword"),
			dbHost:     utils.GetYaml("DBHost"),
			dbName:     utils.GetYaml("DBName")}}
}
