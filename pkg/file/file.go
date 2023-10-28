package file

import (
	"archive/zip"
	"bytes"
	"edetector_go/internal/C_AES"
	"edetector_go/internal/packet"
	"edetector_go/pkg/logger"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CheckDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			logger.Panic("Error creating working dir: " + err.Error())
			panic(err)
		}
		logger.Info("Create dir: " + path)
	}
}

func GetOldestFile(dir string, extension string) (string, string) {
	logCount := 0
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
			logger.Error("Error getting oldest file: " + err.Error())
			time.Sleep(30 * time.Second)
			continue
		}
		if oldestFile == "" {
			logCount += 1
			if logCount == 10 {
				logCount = 0
				logger.Debug("No " + extension + " file to parse")
			}
			time.Sleep(10 * time.Second)
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

func WriteFile(path string, p packet.Packet) error {
	dp := packet.CheckIsData(p)
	decrypt_buf := bytes.Repeat([]byte{0}, len(dp.Raw_data))
	C_AES.Decryptbuffer(dp.Raw_data, len(dp.Raw_data), decrypt_buf)
	decrypt_buf = decrypt_buf[100:]
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(0, 2)
	if err != nil {
		return err
	}
	_, err = file.Write(decrypt_buf)
	if err != nil {
		return err
	}
	return nil
}

func GetFileSize(path string) (int, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	fileLen := fileInfo.Size()
	return int(fileLen), nil
}

func TruncateFile(path string, realLen int) error {
	fileLen, err := GetFileSize(path)
	if err != nil {
		return err
	}
	if int(fileLen) < realLen {
		err = errors.New("incomplete data " + fmt.Sprint(fileLen))
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
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

func UnzipFile(zipPath string, dstPath string, size int) error {
	// truncate data
	err := TruncateFile(zipPath, size)
	if err != nil {
		return err
	}
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
