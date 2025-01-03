package file

import (
	"os"
	"time"

	"github.com/ramil063/gometrics/internal/logger"
)

func retryOpenFile(name string, flag int, perm os.FileMode, tries []int) (*os.File, error) {
	var file *os.File
	var err error
	for try := 0; try < len(tries); try++ {
		time.Sleep(time.Duration(tries[try]) * time.Second)
		file, err = os.OpenFile(name, flag, perm)
		if err == nil {
			break
		}
	}
	return file, err
}

func ClearFileContent(filePath string) error {
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		logger.WriteDebugLog("no file for clear metrics", "ClearFileContent")
		return nil
	}

	err = os.Truncate(filePath, 0)
	return err
}
