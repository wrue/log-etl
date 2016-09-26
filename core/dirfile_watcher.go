package core

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"log-etl/core/collection"
	"os"
	"strings"
	"sync"
	"time"
)

type DirChangeHandler interface {
	fileDeleted(file string)
	fileCreated(file string)
}

var previous *collection.Set

func init() {
	previous = collection.NewSet()
}

func NewDirFileWatcher(dir string, extensions []string, checkPeriod int64) *DirFileWatcher {
	l := list.New()
	return &DirFileWatcher{LogDir: dir, Extensions: extensions,
		CheckPeriod: checkPeriod, Done: false, ListDirChangeHandler: l}
}

type DirFileWatcher struct {
	LogDir               string   //监控目录
	Extensions           []string //监控文件扩展名
	CheckPeriod          int64    //检查周期
	Done                 bool
	ListDirChangeHandler *list.List
	sync.Mutex
}

func (this *DirFileWatcher) AddHandler(hander DirChangeHandler) {
	this.ListDirChangeHandler.PushBack(hander)
}

func ListDir(dirPth string, extensions []string) (files *collection.Set, err error) {
	files = collection.NewSet()
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		for _, extension := range extensions {
			suffix := strings.ToUpper(extension)                       //忽略后缀匹配的大小写
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
				files.Add(dirPth + PthSep + fi.Name())
			}
		}

	}
	return files, nil
}

func (this *DirFileWatcher) check() {
	this.Lock()
	defer this.Unlock()

	files, _ := ListDir(this.LogDir, this.Extensions)
	if files.Len() == 0 {
		fmt.Println("do check ... %s", files.List())
		for _, file := range previous.List() {
			file_str := file.(string)
			this.fireDeletedFile(file_str)
		}
		return
	}

	newFiles := files.Copy()
	addFiles := files.Copy()

	addFiles.RemoveAll(previous)
	for _, af := range addFiles.List() {
		addfile := af.(string)
		this.fireCreatedFile(addfile)
	}
	removedFiles := previous.Copy()
	removedFiles.RemoveAll(newFiles)

	for _, rf := range removedFiles.List() {
		removedFile := rf.(string)
		this.fireDeletedFile(removedFile)
	}
	previous = newFiles

}

func (this *DirFileWatcher) fireDeletedFile(file string) {
	for e := this.ListDirChangeHandler.Front(); e != nil; e = e.Next() {
		//		if s, ok := (e.Value).(DirChangeHandler); ok {
		//			s.fileDeleted(file)
		//		}
		s := e.Value.(DirChangeHandler)
		s.fileDeleted(file)
	}
}

func (this *DirFileWatcher) fireCreatedFile(file string) {
	for e := this.ListDirChangeHandler.Front(); e != nil; e = e.Next() {
		//		if s, ok := (e.Value).(DirChangeHandler); ok {
		//			s.fileDeleted(file)
		//		}
		s := e.Value.(DirChangeHandler)
		s.fileCreated(file)
	}
}

func (this *DirFileWatcher) loop() {
	for {
		if !this.Done {
			this.check()
			time.Sleep(time.Second * 10)
		}
	}
}

func (this *DirFileWatcher) Start() {
	go this.loop()
}
