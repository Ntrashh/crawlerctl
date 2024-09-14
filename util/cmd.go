package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func ExecCmd(command string, args ...string) (string, error) {
	// 创建命令
	cmd := exec.Command(command, args...)

	// 捕获标准输出和错误输出
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	// 运行命令
	err := cmd.Run()
	// 合并输出
	output := outBuffer.String() + errBuffer.String()

	if err != nil {
		return output, err
	}
	return output, nil
}

func EnsureDir(fileName string) error {
	dir := filepath.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("无法创建目录: %v", err)
		}
	}
	return nil
}
func DownloadFile(url string, outputPath string) error {
	// 确保目录存在
	if err := EnsureDir(outputPath); err != nil {
		return err
	}

	// 创建文件
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("无法创建文件: %v", err)
	}
	defer out.Close()

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置 User-Agent，防止某些网站过滤掉默认的请求
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Go-http-client/1.1)")

	// 发送 HTTP GET 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("下载文件出错: %v", err)
	}
	defer resp.Body.Close()
	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 将响应的主体内容写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件时出错: %v", err)
	}
	return nil
}
