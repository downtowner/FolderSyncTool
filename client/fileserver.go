package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
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

//FileServer 文件客户端主服务
type FileServer struct {

	//connect to server
	client *unet.UltralSocket

	//for send heart
	timer *tools.Timer

	//write dir
	dir string

	//parse fileinfo
	filesInfo []DirInfo

	//tmp dir for recv file
	tmpDir string

	//锁定全局事件,如果正在更新文件，那么下一个同步更新请求会被忽略
	lock *sync.Mutex

	syncing bool

	//当前操作的文件
	file *tools.File
}

//NewFileServer 创建文件拉取服务
func NewFileServer() *FileServer {
	p := &FileServer{}
	p.init()
	return p
}

//init 初始化
func (f *FileServer) init() {
	f.timer = tools.NewTimer()
	f.tmpDir = "tmp-CCE7795E10D3"
	f.syncing = false
	f.lock = &sync.Mutex{}
}

//SetWriteDir 设置本地目录
func (f *FileServer) SetWriteDir(dir string) {

	f.dir = tools.CorrectDir(dir)
}

//Run 运行
func (f *FileServer) Run(netaddress string) error {

	if err := f.checkTmpDir(); nil != err {
		return err
	}

	f.client = unet.NewUSocket(netaddress, false)

	err := f.client.Connect(func(mark int8, id int32, cmd string, data []byte) error {

		if mark == int8(2) {
			f.onMessage(cmd, data)
		}
		return nil
	})

	if nil != err {
		log.Println(netaddress, ",连接错误...,err: ", err)
		return err
	}

	f.timer.SetTimer(time.Second*3, func() bool {
		//This usage is not safe
		f.client.SendCmdMessage("heart", nil)
		return true
	})

	log.Println(netaddress, ",连接成功...")

	f.client.Wait()
	return err
}

func (f *FileServer) onMessage(cmd string, data []byte) error {

	switch cmd {
	//文件同步信息
	case "SyncFI":
		return f.onSyncFile(data)
		//开始下载
	case "StartDo":
		return f.onStartDownload(data)
		//下载
	case "DownFi":
		return f.onDownloadFile(data)
		//下载完毕
	case "EndDo":
		return f.onEndDownload()
	}

	return nil
}

func (f *FileServer) checkLocalFile(di *DirInfo) {
	//拼接 文件在本地的路径
	dirPath := tools.CorrectDir(f.dir + di.Dir)
	if !tools.IsExists(dirPath) {
		os.Mkdir(dirPath, 0666)
	}

	//判断文件夹内的文件是否存在
	for _, v := range di.Fis {
		fi := dirPath + v.Name
		if !tools.IsExists(fi) {
			f.requestDownload(di.Dir, v.Name)
			continue
		}

		md5, err := tools.Md5(fi)
		if nil != err {
			continue
		}

		if md5 != v.Md5 {
			f.requestDownload(di.Dir, v.Name)
		}
	}
}

func (f *FileServer) checkTmpDir() error {

	//创建tmp文件夹,防止名字和同步文件夹冲突，名字定为tmp-CCE7795E10D3,同步的文件先放到tmp中，然后覆盖老文件
	if !tools.IsExists(f.dir + f.tmpDir) {
		if err := os.Mkdir(f.dir+f.tmpDir, 0666); nil != err {
			log.Println("创建tmp文件夹失败,err:", err)
			return err
		}
	}

	return nil
}

//download的思想是:先把文件存储在tmp文件夹，接收完毕后复制到目标文件
func (f *FileServer) requestDownload(dir string, filename string) error {
	dlf := RequestDownFi{}
	dlf.Dir = dir
	dlf.FileName = filename

	data, err := json.Marshal(&dlf)
	if nil != err {
		return err
	}

	f.client.SendCmdMessage("ReqDown", data)

	return nil
}

//是否正在同步
func (f *FileServer) isSyncing() bool {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.syncing
}

//
func (f *FileServer) openSync() {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.syncing = true
}

func (f *FileServer) closeSync() {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.syncing = false
}

func (f *FileServer) onDownloadFile(data []byte) error {
	var downfile DownloadFile
	err := json.Unmarshal(data, &downfile)
	if nil != err {
		log.Println("onDownloadFile parse err:", err)
		return nil
	}

	f.file.Write([]byte(downfile.data))

	return nil
}

func (f *FileServer) onStartDownload(data []byte) error {
	var startinfo StartDownloadFi
	err := json.Unmarshal(data, &startinfo)
	if nil != err {
		log.Println("onStartDownload parse err:", err)
		return nil
	}

	//先在tmp目录中组装文件路径
	name := tools.CorrectDir(f.dir+f.tmpDir) + startinfo.FileName
	f.file.OpenWriteFile(name)

	return nil
}

func (f *FileServer) onEndDownload() error {
	f.file.Close()
	return nil
}

func (f *FileServer) onSyncFile(data []byte) error {
	err := json.Unmarshal(data, &f.filesInfo)
	if nil != err {
		log.Println("msg parse err, err:", err)
		return nil
	}

	log.Println("收到消息: SyncFInfo", f.filesInfo[0].Dir)

	for _, v := range f.filesInfo {
		f.checkLocalFile(&v)
	}

	return nil
}
