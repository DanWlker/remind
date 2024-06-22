package helper

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/DanWlker/remind/constant"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetDataFile() string {

	dataFile := strings.TrimSpace(viper.GetString(constant.DATA_FILE_KEY))
	if dataFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			cobra.CheckErr(err)
		}
		dataFile = home + constant.DEFAULT_DATA_PATH_AFTER_HOME + string(os.PathSeparator) + "tempdata.yaml"
	}

	_, errStat := os.Stat(dataFile)
	if errors.Is(errStat, os.ErrNotExist) {
		err := os.MkdirAll(filepath.Dir(dataFile), 0770)
		if err != nil {
			debug.PrintStack()
			log.Println(err)
		}
		if _, errFileCreate := os.Create(dataFile); errFileCreate != nil {
			debug.PrintStack()
			log.Println(err)
		}
	} else if errStat != nil {
		debug.PrintStack()
		log.Println(errStat)
	}

	return dataFile
}
