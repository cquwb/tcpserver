package main

import (
	"fmt"
	"tcp"
	//"time"
	//"net"
)

var g_commandM = tcp.NewCommand()



var g_sess = NewReceiveSession()

func main() {	
	config := &tcp.Config{
		Address:"127.0.0.1:8080",
	}
	g_commandM.Register(2, SS2CS_RegisterServer)
	tcp.BeginClient(g_sess, config)
	/*if !g_sess.WaitRegisterServer() {
		fmt.Printf("register server error")
		return
	}*/
	g_sess.Run()
}

func SS2CS_RegisterServer(s tcp.Sessioner, data []byte) bool {
	if _,ok := s.(*ReceveSession); ok {
		fmt.Printf("get reve data %s \n", data)
	} else {
		return false
	}
	return true
}
