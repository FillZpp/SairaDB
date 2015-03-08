package common

import (
	"os"
)

func IsDirExist(dir string) bool {
	fl, err := os.Stat(dir)

	if err != nil {
		return os.IsExist(err)
	}

	return fl.IsDir()
}

