package main

import (
	"forward/proxy"
)

func main() {
	port := "6666"
	println("启动完成，监听", port, "端口中")
	proxy.Listener(port)
}
