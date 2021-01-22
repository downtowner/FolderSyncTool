package main

func main() {
	fs := NewFileServer()
	fs.SetWriteDir("E:\\local")
	fs.Run("127.0.0.1:10086")
}
