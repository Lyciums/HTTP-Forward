package proxy

import (
	"regexp"
)

var (
	MatchHostPortRegexp = regexp.MustCompile(`(://)?([^:/]+)(:\d+)?`)
)

// HostParseToAddr 从 url 当中解析出域名跟端口
// input : https://httpbin.org
// output: httpbin.org:443
//
// input : http://httpbin.org:8080
// output: httpbin.org:8080
//
// input : httpbin.org:8888
// output: httpbin.org:8888
func HostParseToAddr(host string) string {
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
