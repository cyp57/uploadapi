package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	oid "github.com/coolbed/mgo-oid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Debug(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil { // mean data not obj
		logrus.Debugln(data)
		return
	} else {
		logrus.Debugln(string(bytes))
	}

}

func Output(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}

func ContainInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func GenerateOid(prefix string) string {
	objId := oid.NewOID()

	return prefix + objId.String()
}

func IsEmptyString(s string) bool {
	return s == ""
}

func GetCurrentTime() (time.Time, error) {

	currentTime := time.Now()

	// thaizone, err := time.LoadLocation("Asia/Bangkok")
	// if err != nil {
	// 	return time.Time{}, err
	// }
	// timeString := currentTime.In(thaizone).Format(time.RFC3339)

	timeString := currentTime.Local().Format(time.RFC3339)

	parsedTime, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

func CreateReqFilter(c *gin.Context, arrStr []string) map[string]interface{} {
	filterMap := make(map[string]interface{})
	for _, key := range arrStr {
		if c.Request.URL.Query().Get(key) != "" {
			filterMap[key] = c.Request.URL.Query().Get(key)
		}
	}

	return filterMap
}

func CreateProjection(require map[string]interface{}) (map[string]interface{}, error) {
	var sliceStr []string
	projection := make(map[string]interface{})
	if require["require"] != nil {
		reqStr := fmt.Sprint(require["require"])
		err := json.Unmarshal([]byte(reqStr), &sliceStr)
		if err != nil {
			return nil, err
		}

		for _, value := range sliceStr {
			projection[value] = 1
		}
	} else {
		projection["_id"] = 0
	}
	return projection, nil
}

func GetTypeFromFile(str string) (string, error) {
	fileType := ""
	// Find the last occurrence of "."
	lastDotIndex := strings.LastIndex(str, ".")
	// Extract the substring from the last dot index to the end
	if lastDotIndex != -1 && lastDotIndex < len(str)-1 {
		fileExtension := str[lastDotIndex+1:]
		fmt.Println(fileExtension)
		fileType = fileExtension
	} else {
		return "", errors.New("No file extension found")
	}
	return fileType, nil
}
