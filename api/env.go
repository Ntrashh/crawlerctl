package api

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/models"
	"github.com/Ntrashh/crawlerctl/services"
	"github.com/Ntrashh/crawlerctl/task"
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
	var err error
	if queryType == "pyenv" {
		pythonVersions, err = services.GetPyenvPythonVersions()

	} else if queryType == "virtual" {
		pythonVersions, err = services.GetVirtualPythonVersions()
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
	versions, err := services.GetPyenvPythonVersions()
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	for _, version := range versions {
		if version["version"] == versionData.Version {
			ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Python版本 %s 已经安装!", versionData.Version))
			return
		}
	}
	flag := false
	//TODO 检测正在安装的任务中是否存在
	task.TaskStore.Range(func(key, value interface{}) bool {
		taskStore, ok := value.(*models.Task)
		if !ok {
			return true // 继续遍历
		}
		paramsMap := taskStore.Params.(map[string]interface{})

		// 访问 map 中的值
		version := paramsMap["version"]

		if version == versionData.Version {
			flag = true
		}
		return true
	})

	if flag {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Python版本 %s 正在安装中!", versionData.Version))
		return
	}

	taskFunc := func(params interface{}) (interface{}, error) {
		p := params.(map[string]interface{})
		version := p["version"].(string)

		// 调用服务层的 InstallPyenvPython 方法
		out, installErr := services.InstallPyenvPython(version)
		return map[string]string{
			"message": out,
			"version": version,
		}, installErr
	}

	// 启动任务并传递参数
	taskID := task.StartTask(taskFunc, map[string]interface{}{
		"version": versionData.Version,
	})

	SuccessResponse(c, taskID)
}

func GetRemotePythonVersionHandler(c *gin.Context) {
	versions, err := services.GetRemotePythonVersion()
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, versions)
}
