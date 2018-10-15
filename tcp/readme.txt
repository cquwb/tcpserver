tcp 服务设计
1.server struct {
   Sessions map[id]*Session
   sync.nutex
}
newServer()
NewSession() interface

2.client struct {
   S	*session
}


3 session struct {
	b broker
}
newSession()
Init(broker) { interface
	
}
Process() //处理消息 {
	g_command.Dispatch(this, data[0:4], data[20:])
}

finish() {
	1.处理完已经放在Broker队列里的消息
}



3.server
tcp.NewServer(srv, config) {
1.listen
2. for {
	1.accept
	2.new broker(srv.NewSession(), config)
	3.broker.handler
}

}

4.client
tcp.NewClient(sess, config) {
1.dial
2.new broker(sess,config)
3.broker.handler
}


3.broker {
	net.Conn
	s sess
	ReadChannel []byte
	WriteChannel []*Message
}


handle conn {
	1.go ReadFromConn (stop)
	2.go WriteToConn (stop)
	3.this.s.init()
}

ReadFromConn 从conn度消息到channel
WriteToSession 消息写到SESSION的channel
WriteToConn 从本地的channnel写消息到connn

stop {
	1.this.conn.close() //关闭网络连接，这之前要先处理些什么
}

handleMsg 处理readchannel中的消息，将Broker 传回给session ，然后外面想怎么用就怎么用
this.s.Init()

1.stop and reconnect
a.主动关闭某一个连接
b.读写出现异常关闭
c.程序结束
d.关闭后重连

broker

2.marshal
a.