package file

import (
	"edetector_go/pkg/logger"
	"errors"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func CheckDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			logger.Error("error creating working dir:", zap.Any("error", err.Error()))
		}
		logger.Info("create dir:", zap.Any("message", path))
	}
}

func GetOldestFile(dir string, extension string) string {
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
		return oldestFile
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
		err = errors.New("incomplete data")
		return err
	}
	err = os.WriteFile(path, data[:realLen], 0644)
	if err != nil {
		return err
	}
	return nil
}
