package tcp

import (
	"fmt"
	"net"
	//"time"
	"io"
)

/**
1.简单的一个gorutine来处理一个链接 变成两个 一个read 一个write



*/

type Message struct {
	Head uint32
	Length uint32
	Data []byte
}


type ClientSessions struct {
	Msg chan []byte
	sessions map[string]*ClientSession
}

type ClientSession struct {
	conn net.Conn
	ReadData chan []byte
	WriteMsg chan *Message
	s Sessioner
}

func NewClientSession(c net.Conn, s Sessioner) *ClientSession {
	return &ClientSession{
		conn:c,
		ReadData:make(chan []byte,1024),
		WriteMsg:make(chan *Message,1024),
		s:s,
		
	}
}

func (this *ClientSessions) AddNewSession(s *ClientSession) {
	addr := s.conn.RemoteAddr().String()
	this.sessions[addr] = s
}

func (this *ClientSessions) Begin() bool {
	return true
}

func BeginServer(conf *Config, srv Server) {
	l, e := net.Listen("tcp", conf.Address)
	if e != nil {
		fmt.Printf("[TCPServer] listen error: %v \n", e)
		panic(e.Error())
		//return 这里出错了直接panic 不用return了
	}
	
	defer l.Close()//这个忘了加
	
	fmt.Printf("begin server \n")
	for {
		conn, err := l.Accept()
		/*
		if err != nil {
			fmt.Printf("[TCPServer] accept error: %v \n", err)
			continue
		}
		*/
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue//不懂为什么这么做?背下来就好了
			}
			fmt.Printf("[TCPServer] accept error: %v \n", err)
			return
		}
		fmt.Printf("accpet a connection \n")
		NewClientSession(conn, srv.NewSession()).handleConnection()
	}
}


func BeginClient(sess Sessioner, addr string) {
	con, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("net dial error %v \n", err)
		return
	}
	NewClientSession(con,sess).handleConnection()
}

func (this *ClientSession) readLoop() {
	for {
		
		data := make([]byte, 4, 4)
		n, err := io.ReadFull(this.conn, data)
		if err != nil {
			fmt.Printf("read data error %v \n", err)
			return
		}
		
		fmt.Printf("read data 1 %v \n", data)
		if n != 4 {
			fmt.Printf("read ms begin error %d \n", n)
			continue
		}
		len1 := ByteToUint32(data)
		if len1 < 8 || len1 > 65536 {
			fmt.Printf("read ms len error %d \n", len1)
			continue
		}
		lenNeed := int(len1)
		fmt.Printf("get data len is %d \n", lenNeed)
		data2 := make([]byte, 0, lenNeed+4)
		data2 = append(data2, data...)
		for lenNeed > 0 {
			//这里为什么要用一个for循环来做呢？
			rLen := 4
			if lenNeed < 4 {
				rLen = lenNeed
			}
			data3 := make([]byte, rLen, rLen)
			n2, err2 := this.conn.Read(data3)
			if  err2 != nil {
				fmt.Printf("read data error %v \n", err2)
				return
			}
			fmt.Printf("read data 2 %v \n", data3)
			if n2 < rLen {
				fmt.Printf("read msg body error %d %d \n",len1, n2)
				break
			}
			lenNeed = lenNeed - n2
			data2 = append(data2, data3[:n2]...)
		}
		if len(data2) > 0 {
			this.handleData(data2)
		}
		
	}
}


func (this *ClientSession) writeLoop() {
	for {
		select {
			case msg := <- this.WriteMsg:
				fmt.Printf("WRITE DATA IS %v \n", msg)
				this.conn.Write(Marshal(msg))
		}
	}
}

func Marshal(msg *Message) []byte {
	data := make([]byte, 0, msg.Length)
	data = append(data, Uint32ToByte(msg.Length)...)
	data = append(data, Uint32ToByte(msg.Head)...)
	data = append(data, msg.Data...)
	return data
}

func Unmarshal(data []byte) *Message {
	msg := &Message{}
	msg.Length = ByteToUint32(data[0:4])
	msg.Head = ByteToUint32(data[4:8])
	msg.Data = data[8:]
	return msg
}

/**
这种写法待优化


*/
func (this *ClientSession) Write(head uint32, data []byte) {
	msg := &Message{
		Head:head,
		Length:uint32(4)+uint32(len(data)),
		Data:data,
	}
	this.WriteMsg <- msg
}

func (this *ClientSession) handleConnection() {
	go this.readLoop()
	go this.writeLoop()
	this.s.Init(this)
}

func ByteToUint32(da []byte) uint32 {
	return uint32(da[0]) << 24  + uint32(da[1]) << 16 + uint32(da[2]) << 8 + uint32(da[3])
}

func Uint32ToByte(n uint32) []byte {
	data := make([]byte, 4, 4)
	data[0] = uint8(n >> 24)
	data[1] = uint8(n >> 16)
	data[2] = uint8(n >> 8)
	data[3] = uint8(n)
	return data
}


func (this *ClientSession) handleData(data []byte) {
	fmt.Printf("GET DATA IS %s \n", data)
	//fmt.Printf("get Data length is %d \n", n)
	//todo handle data
	this.ReadData <- data
}
