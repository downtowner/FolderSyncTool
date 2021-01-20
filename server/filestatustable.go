package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
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
	filesInfo map[string][]FileInfo //目录：文件集

	//for lock the filesinfo
	lock *sync.Mutex
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
	f.filesInfo = make(map[string][]FileInfo)
	f.lock.Unlock()

	runtime.GC()
}

func (f *fileStatusTable) Init() {
	f.filesInfo = make(map[string][]FileInfo)
	f.lock = &sync.Mutex{}
}

//WatchDir 监听文件夹
func (f *fileStatusTable) WatchDir(dir string) bool {

	if !f.isDir(dir) {
		log.Println("目标路径: ", dir, "不是文件夹!")
		return false
	}

	dir = f.correctDir(dir)
	f.initfilesInfo(dir)

	return true
}

//WatchFile 监听文件
func (f *fileStatusTable) WatchFile(file string) {
	//是否文件
	if !f.isFile(file) {
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
	md5, err := f.md5(file)
	if nil != err {
		log.Println("获取文件: ", file, "失败, err:", err)
		return
	}

	f.lock.Lock()
	defer f.lock.Unlock()
	f.filesInfo[dirname] = []FileInfo{FileInfo{Name: filename, Md5: md5}}

}

func (f *fileStatusTable) correctDir(dir string) string {
	//判断是否加分隔符
	if "windows" == runtime.GOOS {
		if string(dir[len(dir)-1]) != "\\" {
			dir = dir + "\\"
		}
	} else {
		if string(dir[len(dir)-1]) != "/" {
			dir = dir + "/"
		}
	}

	return dir
}

func (f *fileStatusTable) initfilesInfo(dir string) {

	if f.isDir(dir) {

		sub := strings.Split(dir, "\\")
		key := sub[len(sub)-2]

		filesInfo, err := f.allfilesInfo(dir)
		if nil != err {
			log.Println("目标路径:", dir, "获取文件失败!")
			return
		}

		fis := []FileInfo{}

		for _, v := range filesInfo {

			md5, err := f.md5(dir + v)
			if err != nil {

				log.Println("文件: ", v, "获取md5失败,err:", err)
				continue
			}

			fis = append(fis, FileInfo{Name: v, Md5: md5})
		}

		f.lock.Lock()
		f.filesInfo[key] = fis
		f.lock.Unlock()

		log.Println("==========================================")
		log.Println("文件夹: ", dir)
		log.Println("共检测文件: ", len(fis), "个")
		for k, v := range fis {
			log.Println("序号: ", k+1)
			log.Println("文件名: ", v.Name)
			log.Println("MD5值: ", v.Md5)
		}
		log.Println("==========================================")
	}
}

func (f *fileStatusTable) md5(file string) (string, error) {
	fi, err := os.Open(file)
	if err != nil {

		return "", err
	}
	defer fi.Close()

	body, err := ioutil.ReadAll(fi)
	if err != nil {
		return "", err
	}

	md5 := fmt.Sprintf("%x", md5.Sum(body))
	runtime.GC()
	return md5, nil
}

func (f *fileStatusTable) allfilesInfo(pathname string) ([]string, error) {

	rd, err := ioutil.ReadDir(pathname)
	if nil != err {
		return nil, err
	}

	subfilesInfo := []string{}

	for _, fi := range rd {
		if !fi.IsDir() {
			subfilesInfo = append(subfilesInfo, fi.Name())
		}
	}
	return subfilesInfo, nil
}

func (f *fileStatusTable) exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func (f *fileStatusTable) isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func (f *fileStatusTable) isFile(path string) bool {
	if !f.exists(path) {
		return false
	}

	return !f.isDir(path)
}
