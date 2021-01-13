package main

import (
	"os"
	"syscall"
	"time"

	"git.vnnox.net/ncp/xframe/functional/unet"
)

type watcherserver struct {
	watch *fileWatcher
}

func (w *watcherserver) Init() {
	w.watch = NewFileWatcher()
}

//初始化监听逻辑
func (w *watcherserver) WatchDir(dirs string) {

	w.watch.Listen(dirs)
}

func (w *watcherserver) Run(netaddress string) {

	app := unet.NewTCPServer()

	app.Shutdown(time.Second*5, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	mgr := &LocalManager{}
	mgr.Init()
	mgr.SetConnecter(NewLocalClient)

	app.Run("tcp", netaddress, mgr)
}
