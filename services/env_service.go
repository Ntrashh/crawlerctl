package services

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Ntrashh/crawlerctl/util"
	"os/exec"
	"regexp"
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

func IsOfficialVersion(version string) bool {
	// 正式版本的正则表达式：匹配类似于 "3.8.10" 或 "3.9.5" 这样的版本号
	officialVersionPattern := `^\d+\.\d+\.\d+$`
	re := regexp.MustCompile(officialVersionPattern)
	return re.MatchString(version)
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

// GetPythonVersions 获取当前所有 python 版本
func GetPythonVersions() ([]string, error) {
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

// GetPyenvPythonVersions 查询 pyenv 已安装的 Python 版本
func GetPyenvPythonVersions(versions []string) ([]map[string]interface{}, error) {
	var items []map[string]interface{}
	for version := range versions {
		path, err := GetPyenvVersionPath(versions[version])
		if err != nil {
			return nil, err
		}

		if !IsOfficialVersion(versions[version]) {
			continue
		}
		item := map[string]interface{}{
			"version":  versions[version],
			"path":     path,
			"isGlobal": IsGlobalVersion(versions[version]),
		}
		items = append(items, item)
	}
	return items, nil
}
func ExtractEnvName(path string) (string, bool) {
	// 定义正则表达式，匹配 "/envs/" 后面的部分
	re := regexp.MustCompile(`/envs/([^/]+)`)
	match := re.FindStringSubmatch(path)

	// 如果匹配成功，返回虚拟环境的名字
	if len(match) > 1 {
		return match[1], true
	}
	return "", false
}

func extractVersion(path string) (string, bool) {
	// 定义正则表达式，匹配 `envs` 前面的部分，形如 `3.10.10`
	re := regexp.MustCompile(`^(\d+\.\d+\.\d+)/envs/`)
	match := re.FindStringSubmatch(path)

	// 如果匹配成功，返回版本号
	if len(match) > 1 {
		return match[1], true
	}
	return "", false
}

func GetVirtualPythonVersions(versions []string) ([]map[string]interface{}, error) {
	var items []map[string]interface{}
	for version := range versions {
		path, err := GetPyenvVersionPath(versions[version])
		if err != nil {
			return nil, err
		}

		if IsOfficialVersion(versions[version]) {
			continue
		}
		name, result := ExtractEnvName(versions[version])
		if !result {
			continue
		}
		envVersion, result := extractVersion(versions[version])
		if !result {
			continue
		}
		item := map[string]interface{}{
			"envName": name,
			"version": envVersion,
			"path":    path,
		}
		items = append(items, item)
	}
	return items, nil

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

func GetRemotePythonVersion() ([]string, error) {
	var officialVersions []string
	out, err := util.ExecCmd("pyenv", "install", "-list")
	if err != nil {
		return nil, err
	}
	remoteVersions := strings.Split(out, "\n")
	for _, remoteVersion := range remoteVersions {
		fmt.Printf("remoteVersion,%s \n", remoteVersion)
		remoteVersion = strings.TrimSpace(remoteVersion)
		if IsOfficialVersion(remoteVersion) {
			officialVersions = append(officialVersions, remoteVersion)
		}
	}
	officialVersions = util.ReverseSlice(officialVersions)
	return officialVersions, nil
}
