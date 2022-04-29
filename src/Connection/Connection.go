package main

import (
	"fmt"
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

	//输出PLC运行状态
	fmt.Println(client.PLCGetStatus())
}
