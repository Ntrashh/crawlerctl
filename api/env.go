package api

import (
	"github.com/Ntrashh/crawlerctl/services"
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
	queryType := c.Query("type")
	var pythonVersions []map[string]interface{}
	versions, err := services.GetPythonVersions()
	if queryType == "pyenv" {
		pythonVersions, err = services.GetPyenvPythonVersions(versions)

	} else if queryType == "virtual" {
		pythonVersions, err = services.GetVirtualPythonVersions(versions)
	} else {
		ErrorResponse(c, http.StatusBadRequest, "Invalid query type")
		return
	}
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	SuccessResponse(c, pythonVersions)
}

func InstallPythonHandler(c *gin.Context) {
	var versionData struct {
		Version string `json:"version"`
	}
	err := c.ShouldBindJSON(&versionData)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := services.InstallPyenvPython(versionData.Version)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, result)
}

func GetRemotePythonVersionHandler(c *gin.Context) {
	versions, err := services.GetRemotePythonVersion()
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, versions)
}
