package main

import (
	"log"
	"os"
	"syscall"
	"time"

	"git.vnnox.net/ncp/xframe/functional/unet"
)

//FileInfo 单个文件的状态
type FileInfo struct {
	Name string `json:"name"` //文件名
	Md5  string `json:"md5"`  //md5值
}

//实现文件同步两种方式:
/*
1.启动服务进行同步,监听文件变化然后同步
2.启动服务进行同步,定时检查变更然后同步，目前采用这种
*/
type watcherserver struct {
	//监听器,暂时没用文件监听器,暂时没用到
	//watch *fileWatcher
	//记录监听的目录和文件
	dirs  []string
	files []string

	//定时器
	timer *Timer
}

//Init 初始化
func (w *watcherserver) Init() {
	//w.watch = NewFileWatcher()
	w.timer = &Timer{}
	w.timer.Init()
}

//WatchDir 监听文件夹
func (w *watcherserver) WatchDir(dir string) {
	if !gFileTable.isDir(dir) {
		return
	}

	w.dirs = append(w.dirs, dir)
}

//WatchFile 监听文件
func (w *watcherserver) WatchFile(file string) {
	if !gFileTable.isFile(file) {
		return
	}

	w.files = append(w.files, file)
}

func (w *watcherserver) check() bool {

	if 0 == len(w.dirs) && 0 == len(w.files) {
		log.Println("没有设置监听目录")
		return false
	}

	w.timer.SetTimer(time.Second*10, func() bool {

		gFileTable.Reset()
		for _, v := range w.dirs {
			gFileTable.WatchDir(v)
		}

		for _, v := range w.files {
			gFileTable.WatchFile(v)
		}

		data, _ := gFileTable.Info()
		log.Println(string(data))
		return true
	})

	return true
}

func (w *watcherserver) close() {
	w.timer.Close(0)
	log.Println("关闭所有定时器...")
}

//Run 运行服务
func (w *watcherserver) Run(netaddress string) {

	if !w.check() {
		return
	}

	app := unet.NewTCPServer()

	app.Shutdown(time.Second*5, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	mgr := &LocalManager{}
	mgr.Init()
	mgr.SetExitCallback(w.close)

	mgr.SetConnecter(NewLocalClient)

	app.Run("tcp", netaddress, mgr)
}
