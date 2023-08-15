package checkdir

import (
	"edetector_go/pkg/logger"
	"os"

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
