package main

import "git.vnnox.net/ncp/xframe/functional/unet"

//LocalManager ...
type LocalManager struct {
	unet.ConnectionMgr

	//服务器退出的回调
	willClose func()
}

//Done 覆盖父类的方法
func (l *LocalManager) Done() {

	l.ConnectionMgr.Done()

	if nil != l.willClose {
		l.willClose()
	}
}

//SetExitCallback 设置服务器退出前需要处理的回调
func (l *LocalManager) SetExitCallback(f func()) {
	l.willClose = f
}
