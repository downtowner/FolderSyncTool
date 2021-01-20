package main

import (
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

//暂时没用到监听文件变化

//FileWatcher 简单的封装以下
type fileWatcher struct {
	//监听器
	watcher *fsnotify.Watcher

	//go程管理
	wg *sync.WaitGroup

	//退出信号
	exit chan struct{}
}

//NewFileWatcher 创建监听器
func NewFileWatcher() *fileWatcher {

	p := fileWatcher{}
	p.wg = &sync.WaitGroup{}
	p.exit = make(chan struct{})

	watcher, err := fsnotify.NewWatcher()
	if nil != err {
		return nil
	}

	p.watcher = watcher

	return &p
}

//Listen 监听目录
func (f *fileWatcher) Listen(path string) error {

	if f.isExist(path) {
		f.watcher.Add(path)
	}

	f.watcher.Add(path)
	return nil
}

//Startup 开始监听
func (f *fileWatcher) Startup() {
	//然后再增量同步
	f.localWatch()
}

//Shutdown 关闭监听
func (f *fileWatcher) Shutdown() {
	f.exit <- struct{}{}
	f.watcher.Close()
}

func (f *fileWatcher) isExist(path string) bool {
	_, err := os.Stat(path)
	if nil == err {
		return true
	}

	if os.IsExist(err) {
		return true
	}

	return false
}

func (f *fileWatcher) localWatch() {
	f.wg.Add(1)
	go func() {
		for {

			exit := false
			select {
			case event := <-f.watcher.Events:

				switch event.Op {
				case fsnotify.Create:

				case fsnotify.Write:

				case fsnotify.Remove:

				case fsnotify.Rename:

				case fsnotify.Chmod:

				default:
					log.Println("unknow event:", event.Op.String())

				}

			case err := <-f.watcher.Errors:
				log.Println("watch event error:", err)

			case <-f.exit:
				exit = true
			}

			if exit {
				break
			}
		}

		f.wg.Done()
	}()
}

//Done 完成
func (f *fileWatcher) Done() {
	f.wg.Wait()
	log.Println()
}
