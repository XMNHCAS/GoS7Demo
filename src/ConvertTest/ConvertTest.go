package main

import (
	"fmt"
)

func main() {
	var data byte = 0x12
	value := ByteToBool(data)

	fmt.Println(value)
}

func ByteToBool(data byte) [8]bool {
	var res [8]bool
	for i := 0; i < 8; i++ {
		res[i] = data&1 == 1
		data = data >> 1
	}
	return res
}
