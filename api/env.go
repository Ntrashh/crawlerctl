package api

import (
	"crawlerctl/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckPyenvInstalledHandler 检查 pyenv 是否安装
func CheckPyenvInstalledHandler(c *gin.Context) {
	pyenvInstalled, err := services.CheckPyenvInstalled()

	if err != nil || !pyenvInstalled {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	pyenvVirtualenvInstalled, err := services.CheckPyenvVirtualenvInstalled()
	if err != nil || !pyenvVirtualenvInstalled {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, "Environment created successfully")
}

func GetPyenvPythonVersionHandler(c *gin.Context) {
	versions, err := services.GetPyenvPythonVersions()
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, versions)
}
