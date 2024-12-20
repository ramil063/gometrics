package file

import (
	"os"
	"time"
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
