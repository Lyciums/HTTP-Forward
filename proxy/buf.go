package proxy

import (
	"bytes"
	"io"
	"sync"
)

var (
	DefaultBufferSize = 2048
	// Buffers 降低 GC 压力
	Buffers = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, DefaultBufferSize))
		},
	}
)

func GetBuffer() *bytes.Buffer {
	return Buffers.Get().(*bytes.Buffer)
}

func CopyBufferSize(dst io.Writer, src io.Reader, bfs int) (ws int64) {
	// 如果写入方有 WriteTo 这个方法
	if wt, ok := src.(io.WriterTo); ok {
		ws, _ = wt.WriteTo(dst)
		return
	}
	// 同样，如果目标有 ReadFrom 这个方法
	if rt, ok := dst.(io.ReaderFrom); ok {
		ws, _ = rt.ReadFrom(src)
		return
	}
	if 1 > bfs {
		bfs = DefaultBufferSize
	}
	var (
		buffer = GetBuffer()
		buf    = buffer.Bytes()[:bfs]
	)
	for {
		n, err := src.Read(buf)
		if bfs > n || err == io.EOF {
			break
		}
		n, err = dst.Write(buf)
		if err != nil {
			break
		}
		ws += int64(n)
	}
	return
}

func ReadBufferOfFlag(buffer *bytes.Buffer, r io.Reader, chunkSize int, flag []byte) (*bytes.Buffer, int) {
	if buffer == nil {
		buffer = GetBuffer()
	}
	i, buf, hf := 0, buffer.Bytes(), flag != nil && len(flag) > 0
	// 设置分片大小为 flag 大小
	if 1 > chunkSize {
		// flag 大小为 0，设置默认大小
		chunkSize = DefaultBufferSize
		if hf {
			chunkSize = len(flag)
		}
	}
	fl := chunkSize
	// 按分片大小读取
	for {
		n, err := r.Read(buf[i : i+chunkSize])
		// 到头了
		if fl > n || err == io.EOF {
			break
		}
		// 下一次存储的分片开始位置
		i += n
		// flag 有值则查找新读取的缓存中是否包含 flag
		if hf && bytes.Index(buf[i-n:], flag) > -1 {
			break
		}
	}
	return buffer, i
}
