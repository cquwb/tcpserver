tcp �������
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
Process() //������Ϣ {
	g_command.Dispatch(this, data[0:4], data[20:])
}

finish() {
	1.�������Ѿ�����Broker���������Ϣ
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

ReadFromConn ��conn����Ϣ��channel
WriteToSession ��Ϣд��SESSION��channel
WriteToConn �ӱ��ص�channnelд��Ϣ��connn

stop {
	1.this.conn.close() //�ر��������ӣ���֮ǰҪ�ȴ���Щʲô
}

handleMsg ����readchannel�е���Ϣ����Broker ���ظ�session ��Ȼ����������ô�þ���ô��
this.s.Init()

1.stop and reconnect
a.�����ر�ĳһ������
b.��д�����쳣�ر�
c.�������
d.�رպ�����

broker

2.marshal
a.



3. head {
	Length uint32
	Cmd uint32
	Sid uint64
	Uid uint64
}

4. message {
	Head *PackHead
	Info interface //һ��protobuf �� struct 
}

5.
message => []byte
marshal(message, head, data) []byte, err {
	EncodePackHead(head, message.head)
	return message.Data, nil
}


6.
unmarshal
[]byte => message


7.
message PackHead {
	Length uint32
	Cmd uint32
	Uid uint64
	Sid uint64
}
//��ph ת��Ϊ[]byte
EncodePackHead(buf []byte, ph *PackHead) {
}

//
DecodePackHead(buf []bute, ph *PackHead) {
}
