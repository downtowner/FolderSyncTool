package main

import (
	"encoding/json"
	"log"
	"runtime"
	"strings"
	"sync"
	"tools"
)

//文件状态表,当文件状态变更时，此表应该也更新
var gFileTable *fileStatusTable

func init() {
	gFileTable = &fileStatusTable{}
	gFileTable.Init()
}

//fileStatusTable 文件状态表
type fileStatusTable struct {
	//保存初始的文件md5值
	filesInfo []DirInfo //目录：文件集

	//for lock the filesinfo
	lock *sync.Mutex

	//save file path
	filepath map[string]string
}

//Info 获取文件列表信息
func (f *fileStatusTable) Info() ([]byte, error) {
	f.lock.Lock()
	data, err := json.Marshal(f.filesInfo)
	f.lock.Unlock()

	if nil != err {
		return nil, err
	}

	return data, nil
}

//清空数据
func (f *fileStatusTable) Reset() {
	f.lock.Lock()
	f.filesInfo = nil
	f.lock.Unlock()

	f.filepath = nil

	runtime.GC()
}

//得到真实路径
func (f *fileStatusTable) FilePath(md5 string) string {
	if f, ok := f.filepath[md5]; ok {
		return f
	}

	return ""
}

func (f *fileStatusTable) Init() {
	f.lock = &sync.Mutex{}

	f.filepath = make(map[string]string)
}

//WatchDir 监听文件夹
func (f *fileStatusTable) WatchDir(dir string) bool {

	if !tools.IsDir(dir) {
		log.Println("目标路径: ", dir, "不是文件夹!")
		return false
	}

	dir = tools.CorrectDir(dir)
	f.initfilesInfo(dir)

	return true
}

//WatchFile 监听文件
func (f *fileStatusTable) WatchFile(file string) {
	//是否文件
	if !tools.IsFile(file) {
		log.Println("目标路径: ", file, "不是文件或者目标不存在!")
		return
	}

	var s []string
	if "windows" == runtime.GOOS {
		s = strings.Split(file, "\\")
	} else {
		s = strings.Split(file, "/")
	}

	//得到文件名
	filename := s[len(s)-1]

	//得到父目录
	dirname := s[len(s)-2]

	//md5值
	md5, err := tools.Md5(file)
	if nil != err {
		log.Println("获取文件: ", file, "失败, err:", err)
		return
	}

	f.lock.Lock()
	defer f.lock.Unlock()
	item := FileInfo{Name: filename, Md5: md5}

	dir := DirInfo{}
	dir.Dir = dirname
	dir.Fis = append(dir.Fis, item)

	f.filesInfo = append(f.filesInfo, dir)
}

func (f *fileStatusTable) initfilesInfo(dir string) {

	if tools.IsDir(dir) {

		sub := strings.Split(dir, "\\")
		//文件夹名称
		key := sub[len(sub)-2]

		filesInfo, err := tools.AllFilesInfo(dir)
		if nil != err {
			log.Println("目标路径:", dir, "获取文件失败!")
			return
		}

		fis := []FileInfo{}

		for _, v := range filesInfo {

			md5, err := tools.Md5(dir + v)
			if err != nil {

				log.Println("文件: ", v, "获取md5失败,err:", err)
				continue
			}

			fis = append(fis, FileInfo{Name: v, Md5: md5})

			//真实路径保存下来
			f.filepath[md5] = dir + v
		}

		var dir DirInfo
		dir.Dir = key
		dir.Fis = fis

		f.lock.Lock()
		f.filesInfo = append(f.filesInfo, dir)
		f.lock.Unlock()

		// log.Println("==========================================")
		// log.Println("文件夹: ", dir)
		// log.Println("共检测文件: ", len(fis), "个")
		// for k, v := range fis {
		// 	log.Println("序号: ", k+1)
		// 	log.Println("文件名: ", v.Name)
		// 	log.Println("MD5值: ", v.Md5)
		// }
		// log.Println("==========================================")
	}
}
