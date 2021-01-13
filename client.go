package main

import (
	"log"

	"git.vnnox.net/ncp/xframe/functional/unet"
)

//LocalClient ...
type LocalClient struct {
	unet.ClientSocket
}

//NewLocalClient ...
func NewLocalClient() unet.Client {
	p := &LocalClient{}
	p.Init()
	return p
}

// ...

//OnCmdMessage handle cmd message
func (l *LocalClient) OnCmdMessage(cmd string, data []byte) error {

	log.Println("Local OnCmdMessage")

	return nil
}

//OnIDMessage handle id message
func (l *LocalClient) OnIDMessage(id int, data []byte) error {

	log.Println("Local OnIDMessage", l.RemoteAddress(), "id: ", id, "data: ", string(data))

	l.SendIDMessage(id, data)

	return nil
}
