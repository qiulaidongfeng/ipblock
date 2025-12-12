package ipblock

import (
	"io"
	"log"
	"os"
	"strings"
	"unsafe"
)

// 特殊的标准错误，处理将错误输出到 [os.Stderr],
// 还会检测是否有TLS握手错误，如果有封禁ip。
// 属于实验性API
var Stderr io.Writer = &w{}

type w struct{ r *Rules }

func (wr w) Write(b []byte) (int, error) {
	s := unsafe.String(unsafe.SliceData(b), len(b))
	//可能的错误
	//TLS handshake error from 152.53.22.30:59290: refuse
	//TLS handshake error from 217.119.139.38:34206: client sent an HTTP request to an HTTPS server
	if strings.Contains(s, "TLS handshake error from ") {
		sep := strings.Split(s, "TLS handshake error from ")
		sep = strings.Split(sep[1], ":")
		wr.r.Add(sep[0])
	}
	return os.Stderr.Write(b)
}

// Log 从日志中自动分析并封禁ip
// 属于实验性API
var Log = log.New(Stderr, "", 0)

// Init 将r传递给Stderr
// 属于实验性API
func Init(r *Rules) {
	Stderr.(*w).r = r
}
