package main

import (
	"encoding/json"
	"log"
	"tools"

	"git.vnnox.net/ncp/xframe/functional/unet"
)

const (
	SendBufLen = 64000
)

//LocalClient ...
type LocalClient struct {
	unet.ClientSocket
}

//NewLocalClient ...
func NewLocalClient() unet.Client {
	p := &LocalClient{}
	p.Init()
	return p
}

//OnCmdMessage handle cmd message
func (l *LocalClient) OnCmdMessage(cmd string, data []byte) error {

	switch cmd {
	case "heart":
		log.Println("收到心跳消息...")
	case "ReqDown":
		l.OnReqDownload(data)
	}

	return nil
}

//OnReqDownload 处理请求下载
func (l *LocalClient) OnReqDownload(data []byte) {
	var rdf RequestDownFi
	err := json.Unmarshal(data, &rdf)
	if nil != err {
		log.Println("OnReqDownload err:", err)
		return
	}

	filepath := gFileTable.FilePath(rdf.Md5)
	if "" == filepath {
		return
	}

	//打开这个文件
	f := tools.NewFile()
	f.OpenReadFile(filepath)
	defer f.Close()
	//开始下载
	l.SendCmdMessage("StartDo", data)

	//发送文件
	status := 0
	for {
		buf := make([]byte, SendBufLen)
		n, err := f.Read(buf)
		if nil != err {
			log.Println("file read err: ", err)
			break
		}
		if n < SendBufLen {
			status = 1
			break
		}

		var downfile DownloadFile
		downfile.Info.Dir = rdf.Dir
		downfile.Info.FileName = rdf.FileName
		downfile.Info.Md5 = rdf.Md5
		downfile.data = string(buf[:n])

		data, err = json.Marshal(&downfile)
		if nil != err {
			l.SendCmdMessage("DownFi", data)
		}
	}

	var fs FileStatus
	fs.Status = status
	data, err = json.Marshal(&fs)
	if nil == err {
		log.Println("Marshal fs err:", err)
		return
	}

	l.SendCmdMessage("EndDo", data)
}

//OnIDMessage handle id message
func (l *LocalClient) OnIDMessage(id int, data []byte) error {

	log.Println("Local OnIDMessage", l.RemoteAddress(), "id: ", id, "data: ", string(data))

	l.SendIDMessage(id, data)

	return nil
}

//OnInitOver ..
func (l *LocalClient) OnInitOver() {
	log.Println("我来了...")
	data, err := gFileTable.Info()
	if nil != err {
		log.Println("获取文件列表信息错误,err:", err)
		return
	}

	l.SendCmdMessage("SyncFI", data)
}
