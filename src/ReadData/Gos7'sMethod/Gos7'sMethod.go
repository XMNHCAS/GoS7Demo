package main

import (
	"fmt"
	"time"

	"github.com/robinson/gos7"
)

type PlcData struct {
	boolValue    bool
	intValue     uint16
	realValue    float32
	stringValue  string
	wstringValue string
}

func main() {
	const (
		ipAddr = "192.168.10.230"
		rack   = 0
		slot   = 1
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

	//DB号
	address := 10
	//起始地址
	start := 0
	//读取字节数
	size := 776
	//读写字节缓存区
	buffer := make([]byte, size)

	//读取字节
	client.AGReadDB(address, start, size, buffer)

	//gos7解析数据类
	var helper gos7.Helper

	//gos7内置方法解析数据
	var data PlcData
	data.boolValue = helper.GetBoolAt(buffer[0], 0)
	helper.GetValueAt(buffer[2:4], 0, &data.intValue)
	data.realValue = helper.GetRealAt(buffer[4:8], 0)
	data.stringValue = helper.GetStringAt(buffer[8:264], 0)
	data.wstringValue = helper.GetWStringAt(buffer[264:], 0)

	//输出数据
	fmt.Println(data)
}
