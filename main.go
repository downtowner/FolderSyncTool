package main

func main() {
	app := &watcherserver{}
	app.Init()
	app.WatchDir("")
	app.Run("10086")
}
