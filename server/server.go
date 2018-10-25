package main

import (
	"tcp"
)

/**
*/


type ServerMsg struct {
	sess *ReceveSession
	id uint32
	data []byte
}

type MainServer struct {
	sessions map[string]*ReceveSession
	readData chan *ServerMsg
}

func NewServer() *MainServer {
	return &MainServer{
		sessions:make(map[string]*ReceveSession),
		readData:make(chan *ServerMsg, 1024),
	}
}

func (this *MainServer) PutMsg(data *ServerMsg) {
	this.readData <- data
}

func (this *MainServer) Close() {
	
}

func (this *MainServer) AddSession(ss *ReceveSession) {
	
}


func (this *MainServer) NewSession() tcp.Sessioner {
	return &ReceveSession{}
}

func (this *MainServer) Run() {
	for {
		select {
			case data := <- this.readData:
				this.Process(data)
				
		}
			
	}
}

func (this *MainServer) Process(data *ServerMsg) {
	g_commandM.Dispatch(data.sess, data.id, data.data)
}



type ReceveSession struct {
	sess *tcp.ClientSession
}

func (this *ReceveSession) Init(sess *tcp.ClientSession) {
	this.sess = sess
	go this.Run()
}

func (this *ReceveSession) Run() {
	for {
		select {
			case data := <- this.sess.ReadData:
				this.Process(data)
				
		}
			
	}
}

func (this *ReceveSession) Process(data []byte) {
	msg := tcp.Unmarshal(data)
	g_server.PutMsg(&ServerMsg{this, msg.Head.Id, msg.Data.([]byte)})
	//g_commandM.Dispatch(data)
}

func (this *ReceveSession) Close() {
	
}

func (this *ReceveSession) Write(head uint32, data []byte) {
	this.sess.Write(head, data)
	//g_commandM.Dispatch(data)
}