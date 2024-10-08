package util

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func SaveFileToTemp(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
		}
	}(src)
	dir, err := os.MkdirTemp("", ".crawlerctl")
	if err != nil {
		return "", err
	}
	tempFilePath := filepath.Join(dir, file.Filename)
	dst, err := os.Create(tempFilePath)
	if err != nil {
		return "", err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {

		}
	}(dst)

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return tempFilePath, nil
}

func UnzipFile(zipFilePath, destDir string) error {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
		}
	}(r)

	var commonPrefix string

	// 获取公共前缀（顶层目录名）
	for _, f := range r.File {
		parts := strings.SplitN(f.Name, "/", 2)
		if len(parts) > 1 {
			commonPrefix = parts[0] + "/"
			break
		}
	}

	for _, f := range r.File {
		// 去除顶层目录名
		filePath := strings.TrimPrefix(f.Name, commonPrefix)

		if !isValidZipPath(filePath, destDir) {
			return fmt.Errorf("非法的文件路径：%s", filePath)
		}

		fpath := filepath.Join(destDir, filePath)

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func(outFile *os.File) {
				err := outFile.Close()
				if err != nil {
					fmt.Println(err.Error())
				}
			}(outFile)

			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer func(rc io.ReadCloser) {
				err := rc.Close()
				if err != nil {
					fmt.Println(err.Error())
				}
			}(rc)

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isValidZipPath(filePath, destDir string) bool {
	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return false
	}

	absFilePath, err := filepath.Abs(filepath.Join(destDir, filePath))
	if err != nil {
		return false
	}

	return filepath.HasPrefix(absFilePath, absDestDir)
}

func Base64Encode(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func Base64Decode(text string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		fmt.Println("解码错误:", err)
		return ""
	}
	return string(decodedBytes)
}

func CreateTempDir() string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return ""
	}
	tempDir, _ := os.MkdirTemp("", newUUID.String())
	return tempDir
}
