package proxy

import (
	"fmt"
	"io"
	"net"
)

func Forwarder(client net.Conn) {
	var (
		// 读取首行
		buf, _ = ReadBufferOfFlag(nil, client, 1024, []byte("\n"))
		// 获取请求方法和地址
		method, host string
		_, _         = fmt.Sscanf(buf.String(), "%s%s", &method, &host)
		addr         = HostParseToAddr(host)
		// 建立连接
		server, err = net.Dial("tcp", addr)
	)
	if err != nil {
		Println(addr, "connect failed")
		// Println(addr, "建立链接失败",err)
		// 链接建立失败，关闭客户端链接
		client.Close()
		return
	}
	Println(addr, "connected")
	// connect 方法与其他方法不同
	if method == "CONNECT" {
		client.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	} else {
		// 发送已经读取的部分，这部分从管道中读取完后就消失了，所以需要从 buffer 中读取
		server.Write(buf.Bytes())
		Buffers.Put(buf)
	}
	// 读 client 报文给 server
	go func() {
		// 客户端把请求报文拷贝给服务端完成后，就可以关闭客户端的链接了
		defer client.Close()
		_ = copyConnBody(server, client)
		// if err = copyConnBody(server, client); err != nil{
		// 	Println(addr, "请求远程服务器失败：", err)
		// }
	}()
	// 读 server 响应给 client
	go func() {
		// 服务端把响应体拷贝给客户端完成后，就可以关闭服务端的链接了
		defer server.Close()
		if err := copyConnBody(client, server); err == nil {
			Println(addr, "finished")
			// Println(addr, "响应客户端失败：", err)
		}
	}()
}

func copyConnBody(dst, src net.Conn) error {
	_, err := dst.(io.ReaderFrom).ReadFrom(src)
	return err
}
