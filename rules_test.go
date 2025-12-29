package ipblock

import "testing"

func TestAttack(t *testing.T) {
	if !r.IsBlock("127.0.0.2") {
		t.Fatal("应该阻止")
	}
	Log.Print("TLS handshake error from 152.53.22.30:59290: refuse")
	if !r.IsBlock("152.53.22.30") {
		t.Fatal("应该阻止")
	}
	Log.Print("TLS handshake error from 217.119.139.38:34206: client sent an HTTP request to an HTTPS server")
	if !r.IsBlock("217.119.139.38") {
		t.Fatal("应该阻止")
	}
}

var r = new(Rules)

func init() {
	r.Init("./ipblock.json", nil)
	r.Add("127.0.0.2", "tls scan")
	Init(r)
}
