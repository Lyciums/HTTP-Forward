package proxy

import (
	"log"
)

func init() {
	// 设置日志格式： 时间 文件行
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var (
	Println = log.Println
	Panicln = log.Panicln
)
