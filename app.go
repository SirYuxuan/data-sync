package main

import (
	"data-sync/support"
)

func main() {
	// 启动前置校验/数据库读取
	go support.CheckStart()
	go support.Monitor()
	select {}
}
