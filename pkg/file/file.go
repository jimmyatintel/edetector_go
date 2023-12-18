package file

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
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
			logger.Error("Error creating working dir: " + err.Error())
		}
		logger.Info("Create dir: " + path)
	}
}

func MoveToParentDir() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parentDir := filepath.Dir(currentDir)
	err = os.Chdir(parentDir)
	if err != nil {
		panic(err)
	}
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func GetOldestFile(dir string, extension string) (string, string, string) {
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
		element := strings.Split(oldestFile, "/")
		info := strings.Split(element[len(element)-1], ".")
		err = MoveFile(oldestFile, oldestFile+".processing")
		if err != nil {
			logger.Error("Error renaming file: " + err.Error())
		}
		return oldestFile + ".processing", info[0], info[1]
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

func WriteFile(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(0, 2)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
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

func DecompressionFile(srcPath string, dstPath string, size int) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()
	var firstByte [1]byte
	_, err = file.Read(firstByte[:])
	if err != nil {
		return err
	}
	if firstByte[0] == 'P' {
		err = UnzipFile(srcPath, dstPath, size)
		if err != nil {
			return err
		}
	} else {
		err = UnzipTarFile(srcPath, dstPath, size)
		if err != nil {
			return err
		}
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

func UnzipTarFile(tarPath string, dstPath string, size int) error {
	// truncate data
	err := TruncateFile(tarPath, size)
	if err != nil {
		return err
	}
	// open the tar.gz file for reading
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()
	// create a gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	// create a tar reader
	tarReader := tar.NewReader(gzipReader)
	// extract the files from the tar archive
outerloop:
	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			break outerloop // End of archive
		case err != nil:
			return err
		case header == nil:
			continue // Skip if the header is nil
		}
		// Extract the file
		destFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer destFile.Close()
		_, err = io.Copy(destFile, tarReader)
		if err != nil {
			return err
		}
	}
	// Remove the original file
	err = os.Remove(tarPath)
	if err != nil {
		return err
	}
	return nil
}
