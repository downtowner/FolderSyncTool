package main

import (
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

//FileWatcher 简单的封装以下
type fileWatcher struct {
	watcher *fsnotify.Watcher
	wg      *sync.WaitGroup
	exit    chan struct{}
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
func (f *fileWatcher) Listen(dir string) error {

	if f.isExist(dir) {
		f.watcher.Add(dir)
	}

	return nil
}

func (f *fileWatcher) Startup() {
	//先全部同步检查
	f.syncFolder()
	//然后再增量同步
	f.localWatch()
}

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

func (f *fileWatcher) syncFolder() {

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

func (f *fileWatcher) Done() {
	f.wg.Wait()
	log.Println()
}
