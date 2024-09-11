package services

import (
	"bytes"
	"crawlerctl/util"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// CheckPyenvInstalled 检查 pyenv 是否安装成功
func CheckPyenvInstalled() (bool, error) {
	_, err := exec.LookPath("pyenv")
	if err != nil {
		return false, errors.New("pyenv is not installed")
	}
	return true, nil
}

func CheckPyenvVirtualenvInstalled() (bool, error) {
	// 使用 `pyenv virtualenv --version` 检测是否安装了 pyenv-virtualenv 插件
	cmd := exec.Command("pyenv", "virtualenv", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return false, errors.New("pyenv-virtualenv is not installed or not working properly")
	}
	// 如果命令成功执行并返回版本号，说明安装成功
	version := out.String()
	if version == "" {
		return false, errors.New("pyenv-virtualenv version could not be determined")
	}

	return true, nil
}

func PyenvRootPath() (string, error) {
	out, err := util.ExecCmd("pyenv", "root")
	if err != nil {
		return "", err
	}
	outStr := strings.Replace(out, "\n", "", -1)
	return outStr, nil
}

// GetPyenvVersionPath 获取当前版本的路径
func GetPyenvVersionPath(version string) (string, error) {
	pyenvRootPath, err := PyenvRootPath()
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/versions/%s", pyenvRootPath, version)
	return path, nil
}

// GetPyenvPythonVersions 查询 pyenv 已安装的 Python 版本
func GetPyenvPythonVersions() ([]string, error) {
	// 执行 `pyenv versions --bare` 命令来获取精简的 Python 版本列表
	cmd := exec.Command("pyenv", "versions", "--bare")
	// 捕获命令输出
	output, err := cmd.Output()

	if err != nil {
		return nil, errors.New("failed to execute pyenv command")
	}
	// 将输出按行分割
	versions := strings.Split(string(output), "\n")
	// 去除每行的空格
	var cleanedVersions []string = []string{}
	for _, version := range versions {
		trimmedVersion := strings.TrimSpace(version)
		if trimmedVersion != "" {
			cleanedVersions = append(cleanedVersions, trimmedVersion)
		}
	}
	return cleanedVersions, nil
}

func IsGlobalVersion(version string) bool {
	cmd := exec.Command("pyenv", "global")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	outStr := strings.Replace(string(out), "\n", "", -1)
	if version == outStr {
		return true
	}
	return false
}

func InstallPyenvPython(version string) (bool, error) {
	pyenvRootPath, err := PyenvRootPath()
	downloadUrl := fmt.Sprintf("https://mirrors.huaweicloud.com/python/%s/Python-%s.tar.xz", version, version)
	outputPath := fmt.Sprintf("%s/cache/Python-%s.tar.xz", pyenvRootPath, version)
	err = util.DownloadFile(downloadUrl, outputPath)
	if err != nil {
		return false, err
	}
	_, err = util.ExecCmd("pyenv", "install", version)
	if err != nil {
		return false, err
	}
	return true, nil
}
