package core

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type etlMain struct {
	HOST_NAME                  string
	backUpPath                 string
	processorNum               int
	processWorkDirPath         string
	putClientTransLogPath      string
	processWrokTransLogDirPath string
	Wathcer                    *dirFileWatcher
}

func NewEtlMain() *etlMain {
	logDirPath := viper.GetString("core.logDirPath")
	extensions := viper.GetStringSlice("core.extensions")
	checkPeriod := viper.GetString("core.checkPeriod")

	dirWatcher := newDirFileWatcher(logDirPath, extensions, checkPeriod)
	etl := &etlMain{Wathcer: dirWatcher}
	etl.Wathcer.AddHandler(etl)
	etl.Wathcer.Start()
	for {
		time.Sleep(1 * time.Second)
	}
	return etl
}

func (e *etlMain) fileDeleted(file string) {
	fmt.Printf("fileDeleted %s \n", file)
}

func (e *etlMain) fileCreated(file string) {
	fmt.Printf("fileCreated %s \n", file)
}
