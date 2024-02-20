package utils

import (
	"os"
	"strconv"
)

func GetEnv(key string) string {
	return getEnvVariable(key)
}

func GetEnvBool(key string) bool {
	return getEnvVariableBool(key)
}

func getEnvVariable(key string) string {
	return os.Getenv(key)
}

func getEnvVariableBool(key string) bool {

	result := os.Getenv(key)
	boolValue, err := strconv.ParseBool(result)
	if err != nil {
		return false
	}
	return boolValue

}
