package main

import (
	"tcp"
)

/**
*/


type ServerMsg struct {
	sess *ReceveSession
	data []byte
}



type ReceveSession struct {
	sess *tcp.ClientSession
}

func NewReceiveSession() *ReceveSession {
	return &ReceveSession{}
}

func (this *ReceveSession) Init(sess *tcp.ClientSession) {
	this.sess = sess
	this.RegisterServer(1)
	//go this.Run()
}

func (this *ReceveSession) Close() {
	
}

func (this *ReceveSession) RegisterServer(id uint32) {
	this.Write(1, []byte("你好"))
}

func (this *ReceveSession) Run() {
	for {
		select {
			case data := <- this.sess.ReadData:
				this.Process(data)	
		}	
	}
}

func (this *ReceveSession) Write(head uint32, data []byte) {
	this.sess.Write(head, data)
	//g_commandM.Dispatch(data)
}

func (this *ReceveSession) Process(data []byte) {
	msg := tcp.Unmarshal(data)
	g_commandM.Dispatch(this, msg.Head.Id, msg.Data.([]byte))
	//g_commandM.Dispatch(data)
}


