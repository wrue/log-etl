package core

import (
	"encoding/json"
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
	processTaskChan            chan ProcessTask
	sinkTaskChan               chan SinkTask
	processWrokTransLog        *TransLog
}

func NewEtlMain() *etlMain {
	logDirPath := viper.GetString("core.logDirPath")
	extensions := viper.GetStringSlice("core.extensions")
	checkPeriod := viper.GetString("core.checkPeriod")
	processTaskChan := make(chan ProcessTask, 1024)
	sinkTaskChan := make(chan SinkTask, 1024)

	processWrokTransLog := NewTransLog("E:\\data1\\applog\\trans\\processWorkTransLog")

	processWorkDirPath := "E:\\data1\\applog\\processWorkDirPath"

	dirWatcher := newDirFileWatcher(logDirPath, extensions, checkPeriod)
	etl := &etlMain{Wathcer: dirWatcher}
	etl.processTaskChan = processTaskChan
	etl.sinkTaskChan = sinkTaskChan
	etl.processWorkDirPath = processWorkDirPath
	etl.processWrokTransLog = processWrokTransLog
	etl.Wathcer.AddHandler(etl)
	etl.Wathcer.Start()
	go etl.ProcessWork()
	go etl.SinkerWork()
	for {
		time.Sleep(1 * time.Second)
	}
	return etl
}

func (e *etlMain) SinkerWork() {
	for st := range e.sinkTaskChan {
		fmt.Println(st)
	}
}

func (e *etlMain) ProcessWork() {
	for pt := range e.processTaskChan {
		path := pt.Path
		p := fileToProcessor.GetProcessor(path)
		if p == nil {
			continue
		}
		p.SetWorkDirPath(e.processWorkDirPath)
		pt.SinkTasks = p.Process(path)
		ptBytes, _ := json.Marshal(pt)
		e.processWrokTransLog.WriteAndFlush(string(ptBytes))
		for _, st := range pt.SinkTasks {
			e.sinkTaskChan <- st
		}
	}
}

func (e *etlMain) fileDeleted(file string) {
	fmt.Printf("fileDeleted %s \n", file)
}

func (e *etlMain) fileCreated(file string) {
	fmt.Printf("fileCreated %s \n", file)
	processTask := ProcessTask{}
	processTask.Path = file
	e.processTaskChan <- processTask
}
