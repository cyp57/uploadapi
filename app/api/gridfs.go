package api

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/cyp57/uploadapi/app/controller"
	"github.com/cyp57/uploadapi/app/response"
	"github.com/cyp57/uploadapi/utils"
	"github.com/gin-gonic/gin"
)

func InitApi() IApi {
	return &api{&GridfsApi{}, &FileserverApi{}}
}

type IApi interface {
	Gridfs() IGridfsApi
	Fileserver() IFileserverApi
}

type api struct {
	gridfs     *GridfsApi
	fileserver *FileserverApi
}

func (a *api) Gridfs() IGridfsApi {
	return a.gridfs
}

type IGridfsApi interface {
	UploadFile(c *gin.Context)
	GetFileById(c *gin.Context)
	GetFile(c *gin.Context)
	DeleteFile(c *gin.Context)
}

type GridfsApi struct{}

var responseHandler = response.Response()
var sizeLimit int64 = 10000000  // 10 mb

func (g *GridfsApi) UploadFile(c *gin.Context) {

	file, FileHeader, err := c.Request.FormFile("file")
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	defer file.Close()

	filename := c.Request.FormValue("filename")

	// if FileHeader.Size > sizeLimit { // if want to validate size
	// 	responseHandler.ErrResponse(c, http.StatusBadRequest, "Input file too large")
	// }

	// get type
	split, err := utils.GetTypeFromFile(FileHeader.Filename)
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	var uploadData controller.GridFs
	content, err := io.ReadAll(file)
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	if filename != "" {
		uploadData.FileName = filename
	}
	uploadData.FileType = split
	uploadData.FileSize = FileHeader.Size
	uploadData.Data = content
	result, err := uploadData.UploadFile()
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		responseHandler.SuccessResponse(c, http.StatusOK, result, "Upload success")
	}
}

func (g *GridfsApi) GetFileById(c *gin.Context) {

	id := c.Param("id")

	result, err := new(controller.GridFs).GetFileDetail(id)
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		responseHandler.SuccessResponse(c, http.StatusOK, result, "Success")
	}

}

func (g *GridfsApi) GetFile(c *gin.Context) {

	filename := c.Param("filename")

	size, content, err := new(controller.GridFs).DownloadFile(filename)
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		if size > 0 {
			fileExt := filepath.Ext(filename)

			contentType := mime.TypeByExtension(fileExt)

			if len(contentType) > 0 {
				c.Header("Content-type", contentType)
			} else {
				c.Header("Content-type", "image/*")
			}
			c.Writer.Write(content)
		}
		// responseHandler.SuccessResponse(c,http.StatusOK,result,"Success")
	}

}
func (g *GridfsApi) DeleteFile(c *gin.Context) {

	metadataId := c.Param("id")
	result, err := new(controller.GridFs).DeleteFile(metadataId)
	if err != nil {
		responseHandler.ErrResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		responseHandler.SuccessResponse(c, http.StatusOK, result, "Success")
	}

}
