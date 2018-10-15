package main

import (
	//"fmt"
	"tcp"
)

type ServerConfig struct {
	ipPort uint32
}

var g_server = NewServer()
var g_commandM = tcp.NewCommand()




func RegisterServer(s tcp.Sessioner, data []byte) bool {
	if ss,ok := s.(*ReceveSession); ok {
		ss.Write(2, []byte("ok"))
		g_server.AddSession(ss)
	} else {
		return false
	}
	return true
}

func main() {
	
	g_commandM.Register(1, RegisterServer)
	
	config := &tcp.Config{
		Address:":8080",
	}
	go tcp.BeginServer(config,g_server)
	g_server.Run()
}







