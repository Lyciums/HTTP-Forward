package proxy

import (
	"net"
)

func Listener(port string) {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		Panicln("监听失败：", err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			Panicln("接受请求异常：", err)
		}
		go Forwarder(conn)
	}
}