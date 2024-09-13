package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ntrashh/crawlerctl/util"
	"net/http"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

type EnvService struct {
}

// NewEnvService 创建一个新的 EnvService 实例
func NewEnvService() *EnvService {
	return &EnvService{
		// 初始化依赖或配置
	}
}

// CheckPyenvInstalled 检查 pyenv 是否安装成功
func (s *EnvService) CheckPyenvInstalled() (bool, error) {
	_, err := exec.LookPath("pyenv")
	if err != nil {
		return false, errors.New("pyenv is not installed")
	}
	return true, nil
}

// CheckPyenvVirtualenvInstalled 检查pyenv-virtualenv是否安装
func (s *EnvService) CheckPyenvVirtualenvInstalled() (bool, error) {
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

// PyenvRootPath 获取pyenv的root路径
func (s *EnvService) PyenvRootPath() (string, error) {
	out, err := util.ExecCmd("pyenv", "root")
	if err != nil {
		return "", err
	}
	outStr := strings.Replace(out, "\n", "", -1)
	return outStr, nil
}

// IsOfficialVersion 判断是否是正式版本
func (s *EnvService) IsOfficialVersion(version string) bool {
	// 正式版本的正则表达式：匹配类似于 "3.8.10" 或 "3.9.5" 这样的版本号
	officialVersionPattern := `^\d+\.\d+\.\d+$`
	re := regexp.MustCompile(officialVersionPattern)
	return re.MatchString(version)
}

// GetPyenvVersionPath 获取当前版本的路径
func (s *EnvService) GetPyenvVersionPath(version string) (string, error) {
	pyenvRootPath, err := s.PyenvRootPath()
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/versions/%s", pyenvRootPath, version)
	return path, nil
}

// GetPythonVersions 获取当前所有 python 版本
func (s *EnvService) GetPythonVersions() ([]string, error) {
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
func (s *EnvService) GetPyenvPythonVersions() ([]map[string]interface{}, error) {
	var items []map[string]interface{}
	versions, err := s.GetPythonVersions()
	if err != nil {
		return nil, err
	}
	for version := range versions {
		path, err := s.GetPyenvVersionPath(versions[version])
		if err != nil {
			return nil, err
		}

		if !s.IsOfficialVersion(versions[version]) {
			continue
		}
		item := map[string]interface{}{
			"version":  versions[version],
			"path":     path,
			"isGlobal": s.IsGlobalVersion(versions[version]),
		}
		items = append(items, item)
	}
	return items, nil
}

// extractEnvName 提取虚拟环境名称
func (s *EnvService) extractEnvName(path string) (string, bool) {
	// 定义正则表达式，匹配 "/envs/" 后面的部分
	re := regexp.MustCompile(`/envs/([^/]+)`)
	match := re.FindStringSubmatch(path)

	// 如果匹配成功，返回虚拟环境的名字
	if len(match) > 1 {
		return match[1], true
	}
	return "", false
}

// extractEnvVersion 提取版本
func (s *EnvService) extractEnvVersion(path string) (string, bool) {
	// 定义正则表达式，匹配 `envs` 前面的部分，形如 `3.10.10`
	re := regexp.MustCompile(`^(\d+\.\d+\.\d+)/envs/`)
	match := re.FindStringSubmatch(path)

	// 如果匹配成功，返回版本号
	if len(match) > 1 {
		return match[1], true
	}
	return "", false
}

// GetVirtualPythonVersions  获取虚拟环境的python列表
func (s *EnvService) GetVirtualPythonVersions() ([]map[string]interface{}, error) {
	versions, err := s.GetPythonVersions()
	if err != nil {
		return nil, err
	}
	var items []map[string]interface{}
	for version := range versions {
		path, err := s.GetPyenvVersionPath(versions[version])
		if err != nil {
			return nil, err
		}

		if s.IsOfficialVersion(versions[version]) {
			continue
		}
		name, result := s.extractEnvName(versions[version])
		if !result {
			continue
		}
		envVersion, result := s.extractEnvVersion(versions[version])
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

// IsGlobalVersion 判断是否是全局版本python
func (s *EnvService) IsGlobalVersion(version string) bool {
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

// InstallPyenvPython  安装指定的python版本
func (s *EnvService) InstallPyenvPython(version string) (string, error) {
	pyenvRootPath, err := s.PyenvRootPath()
	downloadUrl := fmt.Sprintf("https://mirrors.huaweicloud.com/python/%s/Python-%s.tar.xz", version, version)
	outputPath := fmt.Sprintf("%s/cache/Python-%s.tar.xz", pyenvRootPath, version)
	err = util.DownloadFile(downloadUrl, outputPath)
	if err != nil {
		return "", err
	}
	out, err := util.ExecCmd("pyenv", "install", version)
	if err != nil {
		return out, err
	}
	return out, nil
}

// GetRemotePythonVersion 获取远程python版本列表
func (s *EnvService) GetRemotePythonVersion() ([]string, error) {
	var officialVersions []string
	out, err := util.ExecCmd("pyenv", "install", "-list")
	if err != nil {
		return nil, err
	}
	remoteVersions := strings.Split(out, "\n")
	for _, remoteVersion := range remoteVersions {
		remoteVersion = strings.TrimSpace(remoteVersion)
		if s.IsOfficialVersion(remoteVersion) {
			officialVersions = append(officialVersions, remoteVersion)
		}
	}
	officialVersions = util.ReverseSlice(officialVersions)
	return officialVersions, nil
}

// SetVersionGlobal 设置指定版本为全局版本
func (s *EnvService) SetVersionGlobal(version string) (string, error) {
	out, err := util.ExecCmd("pyenv", "global", version)
	if err != nil {
		return out, err
	}
	return out, nil
}

// DeletePythonVersion 删除制定的python版本
func (s *EnvService) DeletePythonVersion(version string) (string, error) {

	cmd := exec.Command("sh", "-c", fmt.Sprintf("yes | pyenv uninstall %s", version))
	output, err := cmd.CombinedOutput() // 获取 stdout 和 stderr
	outStr := strings.TrimSpace(string(output))
	if err != nil {
		return outStr, err
	}
	return outStr, nil
}

// DeleteVirtualenv 删除虚拟环境
func (s *EnvService) DeleteVirtualenv(envName string) (bool, error) {

	cmd := exec.Command("sh", "-c", fmt.Sprintf("yes | pyenv virtualenv-delete %s", envName))
	_, err := cmd.CombinedOutput() // 获取 stdout 和 stderr
	if err != nil {
		fmt.Println("Error executing command:", err)
		return false, err
	}
	return true, nil
}

// CreateVirtualenv 创建虚拟环境
func (s *EnvService) CreateVirtualenv(envName, version string) (string, error) {
	out, err := util.ExecCmd("pyenv", "virtualenv", version, envName)
	if err != nil {
		return out, err
	}
	return out, nil
}

// GetVirtualenv 获取指定名称的虚拟环境
func (s *EnvService) GetVirtualenv(envName string) (map[string]interface{}, error) {
	versions, err := s.GetVirtualPythonVersions()
	if err != nil {
		return nil, err
	}
	for _, version := range versions {
		if version["envName"] != envName {
			continue
		}
		return version, nil
	}
	return nil, nil
}

// GetVirtualenvPipPackage 获取虚拟环境已经安装的包
func (s *EnvService) GetVirtualenvPipPackage(path string) ([]map[string]string, error) {
	out, err := util.ExecCmd(fmt.Sprintf("%s/bin/pip", path), "freeze")
	if err != nil {
		return nil, err
	}
	var packages []map[string]string

	// 分割输入字符串为行
	lines := strings.Split(out, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 分割每一行，使用 '==' 作为分隔符
		parts := strings.Split(line, "==")
		if len(parts) == 2 {
			name := parts[0]
			version := parts[1]

			pkg := map[string]string{
				"name":    name,
				"version": version,
			}
			packages = append(packages, pkg)
		}
	}
	return packages, nil
}

func (s *EnvService) GetPackageVersions(packageName string) ([]string, error) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", packageName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("无法发送请求：%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("无法获取包信息，状态码：%d", resp.StatusCode)
	}

	var result struct {
		Releases map[string]interface{} `json:"releases"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败：%v", err)
	}

	var versions []string
	for version := range result.Releases {
		versions = append(versions, version)
	}

	// 对版本号进行排序
	sort.Strings(versions)
	fmt.Println(versions)
	return versions, nil
}
