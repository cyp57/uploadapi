package utils

import (
	"time"

	"github.com/spf13/viper"
)

var ViperYaml *viper.Viper
var ViperJson *viper.Viper

func SetViperYaml(v *viper.Viper) {
	ViperYaml = v
}

func SetViperJson(v *viper.Viper) {
	ViperJson = v
}

func GetViperJson() *viper.Viper {
	return ViperJson
}

func GetYaml(key string) string {
	return getViperVariable(key)
}
func GetYamlDuration(key string) time.Duration {
	return getViperVariableDurationSec(key)
}
func GetYamlInt(key string) int {
	return getViperVariableInt(key)
}
func GetYamlBool(key string) bool {
	return getViperVariableBool(key)
}

func getViperVariable(key string) string {
	value := ViperYaml.GetString(key)

	return value
}

func getViperVariableDurationSec(key string) time.Duration {
	value := ViperYaml.GetDuration(key)

	return value * time.Second
}

func getViperVariableInt(key string) int {
	value := ViperYaml.GetInt(key)

	return value
}

func getViperVariableBool(key string) bool {
	value := ViperYaml.GetBool(key)

	return value
}


