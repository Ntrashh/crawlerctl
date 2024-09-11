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
	versions, err := services.GetPyenvPythonVersions()
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var items []map[string]interface{}
	for version := range versions {
		path, err := services.GetPyenvVersionPath(versions[version])
		if err != nil {
			ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		item := map[string]interface{}{
			"version":  versions[version],
			"path":     path,
			"isGlobal": services.IsGlobalVersion(versions[version]),
		}
		items = append(items, item)
	}
	SuccessResponse(c, items)
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
