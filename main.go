package main

import (
	"forward/proxy"
)

func main() {
	println("启动完成，监听链接中")
	proxy.Proxy("6666")
}
