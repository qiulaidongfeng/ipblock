# ipblock

提供持久化封禁ip的库

将被封禁ip保存在指定路径的json文件

使用示例

```go

package ipblock

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/qiulaidongfeng/ipblock"
)

func Main() {
	r := &Rules{}
	err := r.Init("./ipblock.json")
	if err != nil && err != io.EOF {
		panic(err)
	}
	c := &tls.Config{
		GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
			ip, _, err := net.SplitHostPort(chi.Conn.RemoteAddr().String())
			if err != nil {
				return nil, err
			}
			if r.IsBlock(ip) {
				return nil, errors.New("block")
			}
			if chi.ServerName != "your.domain" {
				r.Add(ip)
				// 设置不同的error便于日志区分执行到那行代码
				return nil, errors.New("refuse")
			}
			// 替换为自己的证书
			return nil, nil
		},
	}
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        // 仅示例，实践中应该在某个中间件进行检测
		if ipblock.MayAttack(req.URL.Path) {
			ip, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				panic(err)
			}
			r.Add(ip)
			return
		}
		fmt.Fprintln(w, "ok")
	})
	s := http.Server{Addr: ":443", TLSConfig: c, Handler: m}
	s.ListenAndServeTLS("", "")
}


```