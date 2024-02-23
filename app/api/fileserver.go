package api

import (
	"net/http"

	"github.com/cyp57/uploadapi/app/controller"
	"github.com/cyp57/uploadapi/utils"
	"github.com/gin-gonic/gin"
)

func (a *api) Fileserver() IFileserverApi {
	return a.fileserver
}

type IFileserverApi interface {
	UploadFile(c *gin.Context)
	GetFileById(c *gin.Context)
	DeleteFile(c *gin.Context)
}

type FileserverApi struct{}

func (f *FileserverApi) UploadFile(c *gin.Context) {

	file, FileHeader, err := c.Request.FormFile("file")
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	defer file.Close()

	filename := c.Request.FormValue("filename")

	if FileHeader.Size > sizeLimit { // if want to validate size
		responseHandler.ErrResponse(c, http.StatusBadRequest, "Input file too large")
	}

	filetype, err := utils.GetTypeFromFile(FileHeader.Filename)
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	var saveData controller.FileServer
	if filename != "" {
		saveData.FileName = filename
	}
	saveData.FileType = filetype
	saveData.FileSize = FileHeader.Size

	result, err := saveData.SaveFile(&file)
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		responseHandler.SuccessResponse(c, http.StatusOK, result, "Upload success")
	}
}

func (f *FileserverApi) GetFileById(c *gin.Context) {

	id := c.Param("id")

	result, err := new(controller.FileServer).GetFileDetail(id)
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		responseHandler.SuccessResponse(c, http.StatusOK, result, "Success")
	}

}

func (f *FileserverApi) DeleteFile(c *gin.Context) {

	id := c.Param("id")
	result, err := new(controller.FileServer).DeleteFile(id)
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		responseHandler.SuccessResponse(c, http.StatusOK, result, "Success")
	}

}
