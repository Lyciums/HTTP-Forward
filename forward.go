package forward

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
)

func StartForwardService(listenPort string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go forward(conn)
	}
}

func forward(conn net.Conn) {
	if conn == nil {
		return
	}
	// 读取 1024 字节的 header，从 http 协议头获取
	var (
		header       [1024]byte
		h            = header[:]
		n, _         = conn.Read(h)
		schemeLine   = string(bytes.SplitN(h, []byte("\n"), 1)[0])
		method, addr string
		_, _         = fmt.Sscanf(schemeLine, "%s%s", &method, &addr)
		uriInfo, _   = url.Parse(addr)
	)
	// 获取连接端口
	if uriInfo.Opaque == "443" {
		addr = uriInfo.Scheme + ":443"
	} else {
		port := uriInfo.Port()
		if port == "" {
			port = "80"
		}
		addr = uriInfo.Host + ":" + port
	}
	log.Println("代理请求：", addr)
	// 连接
	newConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	// 转发
	if method == "CONNECT" {
		conn.Write([]byte("HTTP/1.1 200 fuck proxy\r\n\r\n"))
	} else {
		newConn.Write(h[:n])
	}
	// 写
	go io.Copy(newConn, conn)
	// 读
	io.Copy(conn, newConn)
}