package core

import (
	"bufio"
	"os"
	"strings"
)

type activityinfoLogProcessor struct {
	appLogProcessor
}

func NewActivityinfoLogProcessor() *activityinfoLogProcessor {
	sinkTasks := make([]SinkTask, 0)
	to_out := make(map[string]*os.File)
	activityinfoLogProcessor := &activityinfoLogProcessor{}
	activityinfoLogProcessor.appLogProcessor.sinkTasks = sinkTasks
	activityinfoLogProcessor.appLogProcessor.to_out = to_out
	return activityinfoLogProcessor
}

func (this *activityinfoLogProcessor) Process(filename string) []SinkTask {
	file, _ := os.OpenFile(filename, os.O_RDONLY, 0666)
	this.base = this.workDirPath
	buf := bufio.NewReader(file)
	for {
		lineb, _, err := buf.ReadLine()
		if err != nil {
			break
		}
		line := string(lineb)
		if line == "" || strings.Trim(line, " ") == "" {
			continue
		} else {
			this.print("com/credit/line"+".txt", line)
		}

	}
	return this.sinkTasks
}
