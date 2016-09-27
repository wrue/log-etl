package core

import (
	"container/list"
	"fmt"
	"log-etl/core/collection"
	"log-etl/core/util"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	Period = 60 * 60 * 1000
)

var fileToProcessor AppLogToProcessor = AppLogToProcessor{}

type Processor interface {
	Process(filename string) []SinkTask
	SetWorkDirPath(workDirPath string)
}

type appLogProcessor struct {
	workDirPath string
	sinkTasks   []SinkTask
	to_out      map[string]*os.File
	base        string
}

func (this *appLogProcessor) print(to, line string) {
	out, ok := this.to_out[to]
	if ok {
		out.WriteString(line + "\n")
	} else {
		path := this.base + string(os.PathSeparator) + strings.Replace(to, "/", "_", -1)
		file, _ := os.Create(path)
		sinTask := SinkTask{DataFilePath: path, DestFilePath: to}
		this.sinkTasks = append(this.sinkTasks, sinTask)
		this.to_out[to] = file
		file.WriteString(line + "\n")
	}
}

func (this *appLogProcessor) SetWorkDirPath(workDirPath string) {
	this.workDirPath = workDirPath
}

type FileToProcessor interface {
	GetProcessor(filename string) Processor
}

type AppLogToProcessor struct {
}

func (AppLogToProcessor) GetProcessor(filename string) Processor {
	return NewActivityinfoLogProcessor()
}

type TransLog struct {
	logDirPath string
	rollpoint  int64
	logDir     *os.File
	logFile    *os.File
	sync.Mutex
}

type SinkTask struct {
	DataFilePath           string `json:dataFile`
	DestFilePath           string `json:destFile`
	DestFileOriginalLength int64  `json:destFileOriLen`
}

type ProcessTask struct {
	Path      string     `json:"path"`
	SinkTasks []SinkTask `json:"sinkTasks"`
}

func (this *TransLog) WriteAndFlush(log string) {
	this.Lock()
	defer this.Unlock()
	now := time.Now().UnixNano() / 1e6

	for now > this.rollpoint {
		this.logFile.Close()
		dt := util.GetDateStr(util.GetTime(), util.YYYYMMDD_HH)
		logFilePath := this.logDirPath + string(os.PathSeparator) + "process.data." + dt
		this.rollpoint = this.rollpoint + Period
		logFile, _ := os.Create(logFilePath)
		this.logFile = logFile
	}
	this.logFile.WriteString(log)
	this.logFile.Sync()

}

func (this *TransLog) ReadHours(hours int) *list.List {
	fileInfo, err := this.logDir.Readdir(0)
	if err != nil {
		panic(fmt.Errorf("ReadHours failed ", err))
	}
	sort.Sort(collection.SortedFileArray(fileInfo))
	collection.SortedFileArray(fileInfo).Reverse()
	lines := list.New()
	for i := 0; i < hours; i++ {
		filename := this.logDirPath + string(os.PathSeparator) + fileInfo[i].Name()
		linesArray := util.ReadToStrArray(filename)
		collection.SortedStringArray(linesArray).Reverse()
		for _, line := range linesArray {
			lines.PushBack(line)
		}
	}
	return lines
}

func NewTransLog(logDirPath string) *TransLog {

	_, err := os.Stat(logDirPath)
	exist := os.IsExist(err)
	if !exist {
		os.Mkdir(logDirPath, os.ModePerm)
	}
	logDir, err := os.Open(logDirPath)
	if err != nil {
		fmt.Errorf("Fatal error when opening %s trans log file: %s\n", err)
		os.Exit(-1)
	}
	now := util.GetTime()
	rollpoint := (now.UnixNano() / 1e6) + Period
	dt := util.GetDateStr(now, util.YYYYMMDD_HH)
	logFilePath := logDirPath + string(os.PathSeparator) + "process.data." + dt
	_, err = os.Stat(logFilePath)
	exist = os.IsExist(err)
	if exist {
		fmt.Errorf("Fatal error when created %s trans log file: %s\n", err)
		os.Exit(-1)
	}
	logFile, err := os.Create(logFilePath)
	tranLog := &TransLog{logDirPath: logDirPath, logDir: logDir, rollpoint: rollpoint, logFile: logFile}
	return tranLog
}
