package codecs

import (
	"os"
	"path/filepath"
)

const (
	DONT_PRESERVE_DIR int = 0 // set as zero - default of int type
	PRESERVE_DIR      int = 1
)

var HOME_DIR = getHomeDir()
var BACKUP_DIR = filepath.Join(HOME_DIR, ".backup")

func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}
