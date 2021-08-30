package proxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
)

var (
	BlockedDomainList   []string
	MatchHostPortRegexp = regexp.MustCompile(`(://)?([^:/]+)(:\d+)?`)
	// 降低 GC 压力
	buffers = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 8192))
		},
	}
)

func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

func AddBlockedDomain(domain string) {
	BlockedDomainList = append(BlockedDomainList, domain)
}

func IsBlocked(domain string) bool {
	for _, d := range BlockedDomainList {
		if strings.Index(domain, d) == 0 {
			return true
		}
	}
	return false
}

func Proxy(port string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Panic("监听失败：", err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic("接受请求异常：", err)
		}
		go forwarder(conn)
	}
}

func readBuffer(buffer *bytes.Buffer, r io.ReadCloser, chunkSize int, flag []byte) (*bytes.Buffer, int) {
	if buffer == nil {
		buffer = GetBuffer()
	}
	i, buf := 0, buffer.Bytes()
	// 设置分片大小为 flag 大小
	if 1 > chunkSize {
		// flag 大小为 0，默认 1024
		chunkSize = 1024
		if flag != nil && len(flag) > 0 {
			chunkSize = len(flag)
		}
	}
	// 按分片大小读取
	for {
		n, err := r.Read(buf[i : i+chunkSize])
		// 到头了
		if err == io.EOF {
			break
		}
		// 下一次存储的分片开始位置
		i += n
		// 新读取的缓存中是否包含目标 flag
		if bytes.Index(buf[i-n:], flag) > -1 {
			break
		}
	}
	return buffer, i
}

func hostParseToAddr(host string) string {
	r := MatchHostPortRegexp.FindAllStringSubmatch(host, -1)
	switch len(r) {
	case 1:
		return r[0][0]
	case 2:
		if r[1][3] == "" {
			return r[1][2] + ":80"
		}
		return r[1][2] + r[1][3]
	default:
		return host
	}
}

func forwarder(client net.Conn) {
	var (
		// 读取首行
		buf, _ = readBuffer(nil, client, 1024, []byte("\n"))
		// 获取请求方法和地址
		method, host string
		_, _         = fmt.Sscanf(buf.String(), "%s%s", &method, &host)
		addr         = hostParseToAddr(host)
	)
	// 建立连接
	server, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("与远程服务器建立链接失败：", err)
		return
	}
	log.Println("proxy", addr, method)
	// connect 方法与其他方法不同
	if method == "CONNECT" {
		client.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	} else {
		// 发送已经读取的部分，这部分从管道中读取完后就消失了，所以需要从 buffer 中读取
		server.Write(buf.Bytes())
	}
	// 写
	go func() {
		defer server.Close()
		if _, err := io.Copy(server, client); err != nil {
			log.Println("将数据回写时候发生错误：", err)
		}
	}()
	// 读
	if _, err := io.Copy(client, server); err != nil {
		log.Println("从客户端读取响应时候发生错误", err)
	}
}
