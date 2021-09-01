package proxy

import (
	"net/url"
)

// HostParseToAddr 从 url 当中解析出域名跟端口
func HostParseToAddr(uri string) string {
	var (
		up, _      = url.Parse(uri)
		host, port = up.Hostname(), up.Port()
	)
	if port != "" {
		return host + ":" + port
	}
	port = "443"
	if up.Scheme == "http" {
		port = "80"
	}
	return host + ":" + port
}
