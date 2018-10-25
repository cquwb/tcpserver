package tcp

import (
	"errors"
)

/**
*/

const (
	PACK_HEAD_LEN = 24
	MSG_SIZE = 4
	MSG_MAX_SIZE = 65536
)


var (
	WriteOverflow = errors.New("tcp_conn: write buffer overflow")
	WirteTypeError = errors.New("tcp_conn: write type error")
)

//传递过来的数据用的是big endian 程序运行在little endian
func DecodeUint32(data []byte) uint32 {
	return (uint32(data[0]) << 24) | (uint32(data[1]) << 16) | (uint32(data[2]) << 8) | uint32(data[3])
}

func DecodeUint64(data []byte) uint64 {
	return (uint64(data[0]) << 56) | (uint64(data[1]) << 48) | (uint64(data[2]) << 40) | 
		(uint64(data[3]) << 32) | (uint64(data[4] << 24)) | (uint64(data[5] << 16)) | (uint64(data[6] << 8)) |
		(uint64(data[7]))
}

func EncodeUint32(n uint32, data []byte) {
	data[0] = byte(n >> 24)
	data[1] = byte(n >> 16)
	data[2] = byte(n >> 8)
	data[3] = byte(n)
}

func EncodeUint64(n uint64, data []byte) {
	data[0] = byte(n >> 56)
	data[1] = byte(n >> 48)
	data[2] = byte(n >> 40)
	data[2] = byte(n >> 32)
	data[2] = byte(n >> 24)
	data[2] = byte(n >> 16)
	data[2] = byte(n >> 8)
	data[3] = byte(n)
}