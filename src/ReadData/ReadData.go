package main

import (
	"encoding/binary"
	"fmt"
	//"strings"
	"time"

	"github.com/robinson/gos7"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
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
	helper.GetValueAt(buffer[0:1], 0, &data.boolValue)
	helper.GetValueAt(buffer[2:4], 0, &data.intValue)
	helper.GetValueAt(buffer[4:8], 0, &data.realValue)
	data.stringValue = helper.GetStringAt(buffer[8:264], 0)
	data.wstringValue = helper.GetWStringAt(buffer[264:], 0)

	//输出数据
	fmt.Println(data)

	//手动解析String
	stringPos := 8
	stringData := string(buffer[stringPos+2 : stringPos+2+int(buffer[stringPos+1])])
	fmt.Println(stringData)

	//手动解析WString
	wstringPos := 264
	endPos := binary.BigEndian.Uint16(buffer[wstringPos+2 : wstringPos+4])
	res, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), buffer[wstringPos+4:wstringPos+4+int(endPos)*2])

	//输出数据
	fmt.Println(string(res))
}
