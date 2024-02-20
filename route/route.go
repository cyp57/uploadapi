package route

import (
	"path"

	"github.com/cyp57/uploadapi/app/api"
	"github.com/cyp57/uploadapi/app/middlewares"
	"github.com/cyp57/uploadapi/config"
	"github.com/easonlin404/limit"
	"github.com/gin-gonic/gin"
)

func InitRoute(appConfig config.IAppConfig) *gin.Engine {

	router := gin.Default()

	httpRequestLimit := appConfig.HttpRequestLimit()

	if httpRequestLimit != 0 {
		router.Use(limit.Limit(httpRequestLimit))
	}

	router.Use(new(middlewares.MiddlewareHandler).CorsMiddleware())
	setRoute(router, appConfig)

	router.Run(":" + appConfig.HTTPPort())
	return router
}

func setRoute(router *gin.Engine, appConfig config.IAppConfig) {

	 api := api.InitApi()
	serviceName := appConfig.ServiceName()
	router.GET(serviceName, root)
	v1 := router.Group(serviceName)
	{
		gridfs := v1.Group(path.Join("gridfs"))
		gridfs.POST("/upload",api.Gridfs().UploadFile)
		gridfs.GET("/url/file/:id",api.Gridfs().GetFileById)
		gridfs.DELETE("/file/:id",api.Gridfs().DeleteFile)
		gridfs.GET("/file/:filename",api.Gridfs().GetFile)
	}
	
}

// for health check
func root(c *gin.Context) {
	c.JSON(200, gin.H{"message": "OK"})
}
