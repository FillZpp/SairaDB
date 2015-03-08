package meta

import (
	"fmt"
	"os"
	"path"

	"common"
	"config"
)

var (
	MetaDir string
	NameSpaces map[string]NameSpace
	DefaultNameSapce NameSpace
	Users map[string]User
)

func Init() {
	MetaDir = path.Join(config.ConfMap["data-dir"], "/master/meta")

	if !common.IsDirExist(MetaDir) {
		err := os.MkdirAll(MetaDir, 0700)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"\nError:\nCan not create meta dir %v:\n",
				MetaDir)
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(3)
		}
	}
}

