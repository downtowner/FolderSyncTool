package main

import (
	"log"
	"os"
	"syscall"
	"time"
	"tools"

	"git.vnnox.net/ncp/xframe/functional/unet"
)

//FileInfo 单个文件的状态
type FileInfo struct {
	Name string `json:"name"` //文件名
	Md5  string `json:"md5"`  //md5值
}

//DirInfo 文件夹信息
type DirInfo struct {
	Dir string     `json:"dirname"` //目录名
	Fis []FileInfo `json:"fis"`     //文件信息
}

//RequestDownFi 请求下载文件
type RequestDownFi struct {
	Dir      string `json:"dir"`      //文件夹
	FileName string `json:"filename"` //文件名
	Md5      string `json:"md5"`      //md5
}

//StartDownloadFi 开始下载文件
type StartDownloadFi RequestDownFi

//EndDownloadFi 下载完毕
type EndDownloadFi RequestDownFi

//DownloadFile 下载文件
type DownloadFile struct {
	Info RequestDownFi `json:"info"` //文件信息
	data string        //用string来存放文件数据,接收到需要转成[]byte
}

//FileStatus 文件下载完毕状态
type FileStatus struct {
	Status int `json:"status"` //0:表示下载完成，1:表示因某些错误需要删除
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
	timer *tools.Timer
}

//Init 初始化
func (w *watcherserver) Init() {
	//w.watch = NewFileWatcher()
	w.timer = tools.NewTimer()

}

//WatchDir 监听文件夹
func (w *watcherserver) WatchDir(dir string) {
	if !tools.IsDir(dir) {
		return
	}

	w.dirs = append(w.dirs, dir)
}

//WatchFile 监听文件
func (w *watcherserver) WatchFile(file string) {
	if !tools.IsFile(file) {
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
