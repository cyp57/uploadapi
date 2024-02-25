package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyp57/uploadapi/app/api"
	"github.com/cyp57/uploadapi/app/response"
	"github.com/cyp57/uploadapi/config"
	"github.com/cyp57/uploadapi/pkg/mongodb"
	"github.com/cyp57/uploadapi/setting"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	PathEnv  = "../config/.env"
	PathIni  = "../config/app.ini"
	PathYaml = "../config"
)

func initTestEnvironment() {
	config := config.LoadConfig(PathEnv, PathYaml)
	setting.InitIni(PathIni)
	mongodb.MongoDbConnect(config.Db())
}

func Test_app_Gridfs_UploadFile(t *testing.T) {

	initTestEnvironment()
	router := gin.Default()
	api := api.InitApi()
	router.POST("/gridfs/upload", api.Gridfs().UploadFile)
	router.DELETE("/gridfs/delete/:id", api.Gridfs().DeleteFile)

	// Prepare a test request with a file
	fileContent := "test file content"
	fileName := "testfile.txt"

	// Create a buffer to store the request body
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the file to the request body
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(fileContent))

	// Close the multipart writer to finalize the request body
	writer.Close()

	// Create a mock request with the constructed body
	req, err := http.NewRequest(http.MethodPost, "/gridfs/upload", &requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// Set the Content-Type header with the boundary parameter
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// Create a ResponseRecorder to record the response
	recorder := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(recorder, req)

	// Assert the HTTP status code and response body
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Upload success")

	if recorder.Body != nil {
		responseData := recorder.Body.Bytes()

		var resObj response.ResponseHandler
		err = json.Unmarshal(responseData, &resObj)
		if err != nil {
			t.Fatal(err)
		}
		fileId := resObj.Data.(map[string]interface{})["id"]

		reqDelete, err := http.NewRequest(http.MethodDelete, "/gridfs/delete/"+fmt.Sprint(fileId), nil)
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, reqDelete)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Success")
	}

}

func Test_app_Fileserver_UploadFile(t *testing.T) {
	initTestEnvironment()

	router := gin.Default()
	api := api.InitApi()
	router.POST("/file/upload", api.Fileserver().UploadFile)
	router.DELETE("/file/delete/:id", api.Fileserver().DeleteFile)

	fileContent := "test file content"
	fileName := "testfile.txt"

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(fileContent))

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "/file/upload", &requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Upload success")

	if recorder.Body != nil {
		responseData := recorder.Body.Bytes()

		var resObj response.ResponseHandler
		err = json.Unmarshal(responseData, &resObj)
		if err != nil {
			t.Fatal(err)
		}
		fileId := resObj.Data.(map[string]interface{})["id"]

		reqDelete, err := http.NewRequest(http.MethodDelete, "/file/delete/"+fmt.Sprint(fileId), nil)
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, reqDelete)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Success")
	}

}
