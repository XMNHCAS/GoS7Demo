package main

import (
	"encoding/binary"
	"fmt"
	"math"
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

	var data PlcData

	//手动解析Bool
	data.boolValue = ByteToBool(buffer[0])[0]

	//手动解析Int
	data.intValue = binary.BigEndian.Uint16(buffer[2:4])

	//手动解析Real
	data.realValue = math.Float32frombits(binary.BigEndian.Uint32(buffer[4:8]))

	//手动解析String
	stringPos := 8
	data.stringValue = string(buffer[stringPos+2 : stringPos+2+int(buffer[stringPos+1])])

	//手动解析WString
	wstringPos := 264
	endPos := binary.BigEndian.Uint16(buffer[wstringPos+2 : wstringPos+4])
	res, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), buffer[wstringPos+4:wstringPos+4+int(endPos)*2])
	data.wstringValue = string(res)

	//输出数据
	fmt.Println(data)

}

//字节转bool数组（大端）
func ByteToBool(data byte) [8]bool {
	var res [8]bool
	for i := 0; i < 8; i++ {
		res[i] = data&1 == 1
		data = data >> 1
	}
	return res
}
