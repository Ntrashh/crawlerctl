package api

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/models"
	"github.com/Ntrashh/crawlerctl/services"
	"github.com/Ntrashh/crawlerctl/task"
	"github.com/Ntrashh/crawlerctl/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EnvHandler struct {
	EnvService *services.EnvService
}

// 在初始化时创建 Handler 实例
func NewEnvHandler(envService *services.EnvService) *EnvHandler {
	return &EnvHandler{
		EnvService: envService,
	}
}

// CheckPyenvInstalledHandler 检查 pyenv 是否安装
func (h *EnvHandler) CheckPyenvInstalledHandler(c *gin.Context) {
	pyenvInstalled, err := h.EnvService.CheckPyenvInstalled()

	if err != nil || !pyenvInstalled {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	pyenvVirtualenvInstalled, err := h.EnvService.CheckPyenvVirtualenvInstalled()
	if err != nil || !pyenvVirtualenvInstalled {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, "Environment created successfully")
}

func (h *EnvHandler) GetPyenvPythonVersionHandler(c *gin.Context) {
	queryType := c.Query("type")
	var pythonVersions []map[string]interface{}
	var err error
	if queryType == "pyenv" {
		pythonVersions, err = h.EnvService.GetPyenvPythonVersions()
	} else if queryType == "virtual" {
		pythonVersions, err = h.EnvService.GetVirtualPythonVersions()
	} else {
		ErrorResponse(c, http.StatusBadRequest, "Invalid query type")
		return
	}
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	pythonVersions = util.ReverseSlice(pythonVersions)
	SuccessResponse(c, pythonVersions)
}

func (h *EnvHandler) InstallPythonHandler(c *gin.Context) {
	var versionData struct {
		Version string `json:"version" binding:"required"`
	}
	err := c.ShouldBindJSON(&versionData)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	versions, err := h.EnvService.GetPyenvPythonVersions()
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
	//检测正在安装的任务中是否存在
	task.AsyncTaskStore.Range(func(key, value interface{}) bool {
		taskStore, ok := value.(*models.Task)
		if !ok {
			return true // 继续遍历
		}
		paramsMap := taskStore.Params.(map[string]interface{})

		// 访问 map 中的值
		version := paramsMap["version"]
		status := paramsMap["status"]
		fmt.Println(version, status)
		if version == versionData.Version && status == "running" {
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
		out, installErr := h.EnvService.InstallPyenvPython(version)
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

func (h *EnvHandler) GetRemotePythonVersionHandler(c *gin.Context) {
	versions, err := h.EnvService.GetRemotePythonVersion()
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, versions)
}

func (h *EnvHandler) SetVersionGlobalHandler(c *gin.Context) {
	var versionData struct {
		Version string `json:"version" binding:"required"`
	}
	err := c.ShouldBindJSON(&versionData)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.EnvService.SetVersionGlobal(versionData.Version)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, out)
		return
	}
	SuccessResponse(c, true)
}

func (h *EnvHandler) DeletePythonVersionHandler(c *gin.Context) {
	var versionData struct {
		Version string `json:"version" binding:"required"`
	}
	err := c.ShouldBindJSON(&versionData)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.EnvService.DeletePythonVersion(versionData.Version)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, out)
		return
	}
	SuccessResponse(c, true)
}

func (h *EnvHandler) CreateVirtualenvHandler(c *gin.Context) {
	var versionData struct {
		EnvName string `json:"env_name" binding:"required"`
		Version string `json:"version" binding:"required"`
	}
	err := c.ShouldBindJSON(&versionData)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.EnvService.CreateVirtualenv(versionData.EnvName, versionData.Version)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, out)
		return
	}
	SuccessResponse(c, true)

}

func (h *EnvHandler) DeleteVirtualenvHandler(c *gin.Context) {
	var versionData struct {
		EnvName string `json:"env_name" binding:"required"`
	}
	err := c.ShouldBindJSON(&versionData)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err = h.EnvService.DeleteVirtualenv(versionData.EnvName)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, true)

}

func (h *EnvHandler) GetVirtualenvByNameHandler(c *gin.Context) {
	envName := c.Query("env_name")
	env, err := h.EnvService.GetVirtualenvByName(envName)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if env == nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("未查询到虚拟环境:%s", envName))
		return
	}
	SuccessResponse(c, env)
}

func (h *EnvHandler) VirtualenvPipInstallPackagesHandler(c *gin.Context) {
	var packages = make([]map[string]interface{}, 0)
	var envData struct {
		EnvPath string `json:"env_path"`
	}
	err := c.ShouldBindJSON(&envData)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	packages, err = h.EnvService.GetVirtualenvPipPackage(envData.EnvPath)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, packages)
}

func (h *EnvHandler) GetPackageVersionsHandler(c *gin.Context) {
	packageName := c.Query("package_name")
	packages, err := h.EnvService.GetPackageVersions(packageName)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, packages)

}

func (h *EnvHandler) UninstallPackageHandler(c *gin.Context) {
	var uninstallRequestData struct {
		PackageName    string `json:"package_name" binding:"required"`
		VirtualenvPath string `json:"virtualenv_path" binding:"required"`
	}
	err := c.ShouldBindJSON(&uninstallRequestData)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.EnvService.UninstallPackage(uninstallRequestData.PackageName, uninstallRequestData.VirtualenvPath)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, true)
}

func (h *EnvHandler) InstallPackageHandler(c *gin.Context) {
	var installRequestData struct {
		PackageName        string `json:"package_name" binding:"required"`
		VirtualenvPath     string `json:"virtualenv_path" binding:"required"`
		PackageVersion     string `json:"package_version" binding:"required"`
		InstallationSource string `json:"installation_source"`
	}
	err := c.ShouldBindJSON(&installRequestData)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.EnvService.InstallPackage(installRequestData.PackageName, installRequestData.VirtualenvPath, installRequestData.PackageVersion, installRequestData.InstallationSource)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, true)
}
