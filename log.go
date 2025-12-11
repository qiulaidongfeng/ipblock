package ipblock

import (
	"io"
	"log"
	"os"
	"strings"
	"unsafe"
)

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

var Log = log.New(Stderr, "", 0)

func Init(r *Rules) {
	Stderr.(*w).r = r
}
