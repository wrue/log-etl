package main

import (
	"log-etl/core"
	"time"
)

func main() {
	etlMain := &core.ETLMain{}
	dirWatcher := core.NewDirFileWatcher("E:\\data1\\applog\\log", []string{".pdf"}, 110)
	dirWatcher.AddHandler(etlMain)
	dirWatcher.Start()
	for {
		time.Sleep(1 * time.Second)
	}

}
