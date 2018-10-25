package tcp

import (
	//"fmt"
)

/**
*/



type PackHead struct {
	Length uint32 //包长
	Id uint32 //协议id
	Sid uint64 //服务器id
	Uid uint64 //用户id
}

func GetInputHead(msg []byte) *PackHead {
	return &PackHead{
		Length:uint32(4+len(msg)),
		Id:DecodeUint32(msg[0:4]),
		Sid:DecodeUint64(msg[4:12]),
		Uid:DecodeUint64(msg[12:20]),
	}
}


func EncodePackHead(buf []byte, ph *PackHead) {
	EncodeUint32(ph.Length, buf)
	EncodeUint32(ph.Id, buf[4:])
	EncodeUint64(ph.Sid, buf[8:])
	EncodeUint64(ph.Uid, buf[16:])
}