package tcp

import (
	"fmt"
	"net"
	//"time"
	"io"
	//"sync"
)

/**
1.简单的一个gorutine来处理一个链接 变成两个 一个read 一个write
*/

type Message struct {
	Head *PackHead
	Data interface{}
}


type ClientSession struct {
	//wg sync.WaitGroup
	conn net.Conn
	ReadData chan []byte
	WriteMsg chan *Message
	s Sessioner
	CloseChan chan struct{}
}

func NewClientSession(c net.Conn, s Sessioner) *ClientSession {
	return &ClientSession{
		conn:c,
		ReadData:make(chan []byte,1024),
		WriteMsg:make(chan *Message,1024),
		s:s,
		CloseChan:make(chan struct{}),	
	}
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


func BeginClient(sess Sessioner, conf *Config) {
	con, err := net.Dial("tcp", conf.Address)
	if err != nil {
		fmt.Printf("net dial error %v \n", err)
		return
	}
	NewClientSession(con, sess).handleConnection()
}

func (this *ClientSession) readLoop() {
	for {
		data := make([]byte, MSG_SIZE, MSG_SIZE)
		n, err := io.ReadFull(this.conn, data)
		if err != nil {
			fmt.Printf("read data error %v \n", err)
			this.Stop()
			goto EXIT
		}
		
		fmt.Printf("read data 1 %v %d \n", data, n)
		if n != MSG_SIZE {
			fmt.Printf("read ms begin error %d \n", n)
			continue
		}
		len1 := DecodeUint32(data)
		if len1 < PACK_HEAD_LEN-MSG_SIZE || len1 > 65536 {
			fmt.Printf("read ms len error %d \n", len1)
			continue
		}
		lenNeed := int(len1)-MSG_SIZE
		fmt.Printf("read need len %d \n", lenNeed)
		
		data12 := make([]byte, lenNeed, lenNeed)
		_, err2 := io.ReadFull(this.conn, data12)
		if err2 != nil {
			fmt.Printf("read data error %v \n", err2)
			this.Stop()
			goto EXIT
		}
		if len(data12) > 0 {
			this.handleData(data12)
		}
		
	}
EXIT:
	//todo anything
}


func (this *ClientSession) writeLoop() {
	writeData := make([]byte, MSG_MAX_SIZE)
	HeadData := make([]byte, PACK_HEAD_LEN)
	DataBuff := make([]byte, MSG_MAX_SIZE-PACK_HEAD_LEN)
	index := 0
	
	for {
		select {
			case msg := <- this.WriteMsg:
				fmt.Printf("Write Data Head is %v \n", msg.Head)
				fmt.Printf("Write Data Data is %v \n", msg.Data)
				n , data, err := Marshal(msg, HeadData, DataBuff)
				if err != nil {
					break
				}
				copy(writeData[0:], HeadData)
				copy(writeData[PACK_HEAD_LEN:], data)
				index += n
				for more:=true;more; {
					select {
						case msg:=<-this.WriteMsg:
							n , data, err := Marshal(msg, HeadData, DataBuff)
							if err != nil {
								more = false
								break
							}
							if index + n > MSG_MAX_SIZE {
								if _, err := this.conn.Write(writeData[0:index]); err != nil {
									this.Stop()
									goto EXIT
								}
								index = 0
								copy(writeData[0:], HeadData)
								copy(writeData[PACK_HEAD_LEN:], data)
								more = false
							} else {
								copy(writeData[index:], HeadData)
								copy(writeData[index+PACK_HEAD_LEN:], HeadData)
								index += n
							}
						case <- this.CloseChan:
							goto EXIT
						default:
							more = false	
					}
				}
	
				if _, err := this.conn.Write(writeData[0:index]); err != nil {
					this.Stop()
					goto EXIT
				}
			case <-this.CloseChan:
				goto EXIT
		}
	}
EXIT:
}

func Marshal(msg *Message, HeadBuff []byte, DataBuff []byte) (int, []byte, error) {
	
	
	
	switch  v:=msg.Data.(type) {
		case  []byte: 
			DataBuff =v
		default:
			fmt.Printf("write data type error \n")
			return 0, nil, WirteTypeError
	}
	
	if len(DataBuff)+PACK_HEAD_LEN > MSG_MAX_SIZE {
		fmt.Printf("write data type over flow \n")
		return 0, nil, WriteOverflow
	}
	
	EncodePackHead(HeadBuff, msg.Head)
	
	return int(msg.Head.Length), DataBuff, nil
}

func Unmarshal(data []byte) *Message {
	
	head := GetInputHead(data)
	msg := &Message{
		head,
		data[20:],
	}
	
	return msg
}



/**
这种写法待优化


*/
func (this *ClientSession) Write(head uint32, data []byte) {
	msg := &Message{
		Head:&PackHead{uint32(PACK_HEAD_LEN+len(data)),head, 0, 0},
		Data:data,
	}
	this.WriteMsg <- msg
}

func (this *ClientSession) handleConnection() {
	//this.wg.Add()
	go this.readLoop()
	//this.wg.Add()
	go this.writeLoop()
	this.s.Init(this)
}




func (this *ClientSession) handleData(data []byte) {
	fmt.Printf("Read Data byte is %v \n", data)
	fmt.Printf("Read Data string is %s \n", data)
	//fmt.Printf("get Data length is %d \n", n)
	//todo handle data
	this.ReadData <- data
}

//todo 
func (this *ClientSession) Stop() {
	
}
