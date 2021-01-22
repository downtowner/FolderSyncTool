package tools

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

//Md5 获得文件md5值
func Md5(file string) (string, error) {
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

//AllFilesInfo 获取当前1级路径下所有文件
func AllFilesInfo(pathname string) ([]string, error) {

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

//IsExists 判断文件是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

//IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	if !IsExists(path) {
		return false
	}

	return !IsDir(path)
}

//CorrectDir 矫正路径
func CorrectDir(dir string) string {
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
