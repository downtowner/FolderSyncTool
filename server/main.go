package main

func main() {
	app := &watcherserver{}
	app.Init()
	app.WatchDir("E:\\watcher")
	app.WatchDir("E:\\dbfile")
	app.WatchFile("E:\\webView.zip")
	app.Run(":10086")
}
