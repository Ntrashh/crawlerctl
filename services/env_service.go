package services

import (
	"bytes"
	"errors"
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
