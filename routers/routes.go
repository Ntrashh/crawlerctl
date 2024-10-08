package routes

import (
	"github.com/Ntrashh/crawlerctl/api"
	"github.com/Ntrashh/crawlerctl/middleware"
	"github.com/Ntrashh/crawlerctl/services"
	"github.com/Ntrashh/crawlerctl/storage"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {

	router.POST("/login", api.LoginHandler)
	// 环境管理相关的路由
	envService := services.NewEnvService()
	envHandler := api.NewEnvHandler(envService)
	envRoutes := router.Group("/envs")
	envRoutes.Use(middleware.AuthMiddleware())
	{
		envRoutes.GET("/check_install", envHandler.CheckPyenvInstalledHandler)
		envRoutes.GET("/get_versions", envHandler.GetPyenvPythonVersionHandler)
		envRoutes.POST("/install", envHandler.InstallPythonHandler)
		envRoutes.GET("/remote_versions", envHandler.GetRemotePythonVersionHandler)
		envRoutes.POST("/set_global", envHandler.SetVersionGlobalHandler)
		envRoutes.POST("/delete_python", envHandler.DeletePythonVersionHandler)
		envRoutes.POST("/create_virtualenv", envHandler.CreateVirtualenvHandler)
		envRoutes.POST("/delete_virtualenv", envHandler.DeleteVirtualenvHandler)
		envRoutes.GET("/get_virtualenv", envHandler.GetVirtualenvByNameHandler)
		envRoutes.POST("/installed_packages", envHandler.VirtualenvPipInstallPackagesHandler)
		envRoutes.GET("/package_versions", envHandler.GetPackageVersionsHandler)
		envRoutes.POST("/uninstall_packages", envHandler.UninstallPackageHandler)
		envRoutes.POST("/install_packages", envHandler.InstallPackageHandler)
		envRoutes.POST("/install_requirements", envHandler.InstallRequirementsHandler)
	}
	projectStorage := storage.NewProjectStorage()
	projectService := services.NewProjectService(projectStorage)
	projectHandler := api.NewProjectHandler(projectService)
	projectRoutes := router.Group("/projects")
	projectRoutes.Use(middleware.AuthMiddleware())
	{
		projectRoutes.GET("/projects", projectHandler.GetAllProjectsHandler)
		projectRoutes.POST("/add_project", projectHandler.AddProjectHandler)
		projectRoutes.DELETE("/:id", projectHandler.DeleteProjectHandler)
		projectRoutes.GET("/projects_by_version", projectHandler.ProjectsByVersionHandler)
		projectRoutes.GET("/get_folders", projectHandler.GetFolderContentsHandler)
		projectRoutes.GET("/read_file", projectHandler.ReadFileHandler)
		projectRoutes.GET("/:id", projectHandler.ProjectByIdHandler)
		projectRoutes.POST("/save_file", projectHandler.SaveFileHandler)
		projectRoutes.POST("/reload_file", projectHandler.ReUploadHandler)
	}
	gitStorage := storage.NewGitStore()
	gitService := services.NewGitService(gitStorage)
	gitHandler := api.NewGitHandler(gitService, projectService)
	gitRoutes := router.Group("/git")
	gitRoutes.Use(middleware.AuthMiddleware())
	{
		gitRoutes.GET("/:id", gitHandler.GitByProjectIdHandler)
		gitRoutes.POST("/create_git", gitHandler.CreateGitConfigHandler)
		gitRoutes.GET("/remote_branches/:id", gitHandler.RemoteBranchesHandler)
		gitRoutes.GET("/remote_commits", gitHandler.RemoteBranchCommitsHandler)
		gitRoutes.POST("/breach_pull", gitHandler.BranchPullHandler)

	}

	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("/task_status", api.GetTaskStatus)
	}

}
