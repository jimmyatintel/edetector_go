package file

import (
	"archive/zip"
	"edetector_go/pkg/logger"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

func CheckDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			logger.Panic("error creating working dir:", zap.Any("error", err.Error()))
			panic(err)
		}
		logger.Info("create dir:", zap.Any("message", path))
	}
}

func GetOldestFile(dir string, extension string) (string, string) {
	for {
		var oldestFile string
		var oldestTime time.Time
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == extension {
				modTime := info.ModTime()
				if oldestTime.IsZero() || modTime.Before(oldestTime) {
					oldestTime = modTime
					oldestFile = path
				}
			}
			return nil
		})
		if err != nil {
			logger.Error("Error getting oldest file:", zap.Any("error", err.Error()))
			time.Sleep(30 * time.Second)
			continue
		}
		if oldestFile == "" {
			logger.Info("No file to parse")
			time.Sleep(30 * time.Second)
			continue
		}
		element := strings.Split(strings.Split(oldestFile, extension)[0], "/")
		agent := strings.Split(element[len(element)-1], "-")[0]
		return oldestFile, agent
	}
}

func CreateFile(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func WriteFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(0, 2)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func TruncateFile(path string, realLen int) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	fileLen := fileInfo.Size()
	if int(fileLen) < realLen {
		err = errors.New("incomplete data " + fmt.Sprint(fileLen))
		return err
	}
	err = os.WriteFile(path, data[:realLen], 0644)
	if err != nil {
		return err
	}
	return nil
}

func MoveFile(srcPath string, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	err = os.Remove(srcPath)
	if err != nil {
		return err
	}
	return nil
}

func UnzipFile(zipPath string, dstPath string) error {
	// open the zip file for reading
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	// extract the files from the zip archive
	for _, file := range reader.File {
		if !file.FileInfo().IsDir() {
			destFile, err := os.Create(dstPath)
			if err != nil {
				return err
			}
			srcFile, err := file.Open()
			if err != nil {
				return err
			}
			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}
			destFile.Close()
			srcFile.Close()
		} else {
			err = errors.New("the zip file contains a directory")
			return err
		}
	}
	reader.Close()
	err = os.Remove(zipPath)
	if err != nil {
		return err
	}
	return nil
}
