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

func ClearDirContent(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	CheckDir(path)
	return nil
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

func CopyFile(srcPath string, dstPath string) error {
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
		err = UnZipFile(srcPath, dstPath, size)
		if err != nil {
			return err
		}
	} else {
		err = UnTarFile(srcPath, dstPath, size)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnZipFile(zipPath string, dstPath string, size int) error {
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

func UnTarFile(tarPath string, dstPath string, size int) error {
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

func DecompressionDir(srcPath string, dstPath string) error {
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
		err = UnZipDir(srcPath, dstPath)
		if err != nil {
			return err
		}
	} else {
		err = UnTarGz(srcPath, dstPath)
		if err != nil {
			return err
		}
	}
	// remove the original file
	err = os.Remove(srcPath)
	if err != nil {
		return err
	}
	return nil
}

func UnZipDir(zipPath string, dstPath string) error {
	// open the zip file for reading
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	// extract the files from the zip archive
	for _, file := range reader.File {
		dstFilePath := filepath.Join(dstPath, file.Name)
		if !file.FileInfo().IsDir() {
			// open the source file
			srcFile, err := file.Open()
			if err != nil {
				return err
			}
			defer srcFile.Close()
			// create the directory for the file
			err = os.MkdirAll(filepath.Dir(dstFilePath), 0755)
			if err != nil {
				return err
			}
			// create the destination file
			dstFile, err := os.Create(dstFilePath)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			// copy the contents of the file
			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				return err
			}
		} else {
			// create directories if it is a directory
			err = os.MkdirAll(dstFilePath, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UnTarGz(tarGzPath string, dstPath string) error {
	// open the tar.gz file for reading
	file, err := os.Open(tarGzPath)
	if err != nil {
		return err
	}
	defer file.Close()
	// create a gzip reader
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()
	// create a tar reader
	tr := tar.NewReader(gzr)
	// extract the files from the tar archive
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // end of archive
		}
		if err != nil {
			return err
		}
		dstFilePath := filepath.Join(dstPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			// create directory if it is a directory
			if err := os.MkdirAll(dstFilePath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// create the directory for the file
			if err := os.MkdirAll(filepath.Dir(dstFilePath), 0755); err != nil {
				return err
			}
			// create the destination file
			dstFile, err := os.Create(dstFilePath)
			if err != nil {
				return err
			}
			defer dstFile.Close()
			// copy the contents of the file
			if _, err := io.Copy(dstFile, tr); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported file type: %v", header.Typeflag)
		}
	}
	return nil
}

func ZipDir(srcPath string, dstPath string) error {
	// create the destination zip file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	// create a zip.Writer
	zw := zip.NewWriter(dstFile)
	defer zw.Close()
	// walk through all files and subdirectories in the source directory
	err = filepath.Walk(srcPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// create a zip file header information
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		// modify the file name in the zip header to be a relative path
		relPath, err := filepath.Rel(srcPath, filePath)
		if err != nil {
			return err
		}
		header.Name = relPath
		// for directories, add a trailing slash to the name
		if info.IsDir() {
			header.Name += "/"
		}
		// write the zip file header information
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		// if it's a file, write the file content to the zip file
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(writer, file); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	// remove the source directory
	err = os.RemoveAll(srcPath)
	if err != nil {
		return err
	}
	return nil
}

func TarDir(srcPath string, dstPath string) error {
	// Create the destination tar file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	// Create a tar.Writer
	tw := tar.NewWriter(dstFile)
	defer tw.Close()
	// Walk through all files and subdirectories in the source directory
	err = filepath.Walk(srcPath, func(filePath string, info os.FileInfo, err error) error {
		// Ignore the source directory itself
		if filePath == srcPath {
			return nil
		}
		if err != nil {
			return err
		}
		// Create tar file header information
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		// Modify the file name in the tar header to be a relative path
		relPath, err := filepath.Rel(srcPath, filePath)
		if err != nil {
			return err
		}
		header.Name = relPath
		// Write the tar file header information
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// If it's a file, write the file content to the tar file
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	// remove the source directory
	err = os.RemoveAll(srcPath)
	return nil
}
