package tools

import "os"

//File 封装对文件的操作
type File struct {
	file *os.File
}

//NewFile 新建文件
func NewFile() *File {
	return &File{}
}

//OpenWriteFile ...
func (f *File) OpenWriteFile(filename string) error {
	var err error
	f.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0666)
	return err
}

//OpenReadFile 只读的方式打开文件
func (f *File) OpenReadFile(filename string) error {
	var err error
	f.file, err = os.Open(filename)
	return err
}

//Write 写文件
func (f *File) Write(data []byte) (int, error) {
	return f.file.Write(data)
}

//Read 读文件
func (f *File) Read(buf []byte) (int, error) {
	return f.file.Read(buf)
}

//Close 关闭文件
func (f *File) Close() {
	f.file.Close()
}
