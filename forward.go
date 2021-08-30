package forward

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
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
		// https 的 host 在 scheme 里
		addr = uriInfo.Scheme + ":443"
	} else if strings.Index(uriInfo.Host, ":") > -1 {
		// 自带端口
		addr = uriInfo.Host
	} else {
		// 默认 80 端口
		addr = uriInfo.Host + ":80"
	}
	log.Println("代理请求：", addr)
	// 连接
	newConn, err := net.DialTimeout("tcp", addr, time.Second*4)
	if err != nil {
		log.Println("拨号失败：", err)
		return
	}
	// 结束后关闭链接
	defer func() {
		conn.Close()
		newConn.Close()
	}()
	// 转发
	if method == "CONNECT" {
		conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	} else {
		newConn.Write(h[:n])
	}
	// 写
	go func() {
		if _, err := io.Copy(newConn, conn); err != nil {
			log.Println("将数据回写时候发生错误：", err)
		}
	}()
	// 读
	if _, err := io.Copy(conn, newConn); err != nil {
		log.Println("从客户端读取响应时候发生错误", err)
	}
}
