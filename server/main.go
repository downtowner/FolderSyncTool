package main

func main() {
	app := &watcherserver{}
	app.Init()
	app.WatchDir("E:\\watcher")
	app.Run(":10086")
}
