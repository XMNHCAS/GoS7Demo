package main

import (
	"time"

	"github.com/robinson/gos7"
)

func main() {
	const (
		ipAddr = "192.168.10.230" //PLC IP
		rack   = 0                // PLC机架号
		slot   = 1                // PLC插槽号
	)
	//PLC tcp连接客户端
	handler := gos7.NewTCPClientHandler(ipAddr, rack, slot)
	//连接及读取超时
	handler.Timeout = 200 * time.Second
	//关闭连接超时
	handler.IdleTimeout = 200 * time.Second
	//打开连接
	handler.Connect()
	//函数退出时关闭连接
	defer handler.Close()

	//获取PLC对象
	client := gos7.NewClient(handler)

	//gos7解析数据类
	var helper gos7.Helper

	//DB号
	dbNo := 10
	//起始地址
	startAdr := 0

	//写入数据的字节二位数组
	buffers := [][]byte{
		make([]byte, 2),
		make([]byte, 2),
		make([]byte, 4),
		make([]byte, 256),
		make([]byte, 512),
	}

	//生成需要写入的变量的数组
	helper.SetBoolAt(buffers[0][0], 0, false)
	helper.SetValueAt(buffers[1], 0, uint16(100))
	helper.SetRealAt(buffers[2], 0, float32(66.6))
	helper.SetStringAt(buffers[3], 0, 254, "Hello Go")
	SetWStringAt(buffers[4], 0, "中文")

	//循环数据，逐个写入
	for _, v := range buffers {
		client.AGWriteDB(dbNo, startAdr, len(v), v)
		startAdr += len(v)
	}
}

//获取WString的报文
func SetWStringAt(buffer []byte, pos int, value string) []byte {
	chars := []rune(value)
	slen := len(chars)
	var maxLen int = 254
	if maxLen < slen {
		maxLen = slen
	}
	var helper gos7.Helper
	helper.SetValueAt(buffer, pos+0, int16(maxLen))
	helper.SetValueAt(buffer, pos+2, int16(slen))
	for i, c := range chars {
		if i >= maxLen {
			return buffer
		}
		helper.SetValueAt(buffer, pos+4+i*2, uint16(c))
	}
	return buffer
}
