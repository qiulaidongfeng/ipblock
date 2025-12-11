package ipblock

import (
	"encoding/json"
	"os"
	"sync"
)

// Rules 表示被阻止ip的规则集合
type Rules struct {
	ips  sync.Map
	path string
	m    sync.Mutex
}

func (r *Rules) Add(ip string) {
	r.ips.Store(ip, struct{}{})
	if err := r.Save(); err != nil {
		panic(err)
	}
}

func (r *Rules) IsBlock(ip string) bool {
	_, ok := r.ips.Load(ip)
	return ok
}

func (r *Rules) Save() error {
	r.m.Lock()
	defer r.m.Unlock()
	fd, err := os.Create(r.path)
	if err != nil {
		return err
	}
	d := json.NewEncoder(fd)
	d.SetIndent("", "\t")
	var m = make(map[string]struct{})
	r.ips.Range(func(key, value any) bool {
		m[key.(string)] = struct{}{}
		return true
	})
	err = d.Encode(&m)
	return err
}

func (r *Rules) Init(path string) error {
	r.path = path
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	d := json.NewDecoder(fd)
	var m = make(map[string]struct{})
	err = d.Decode(&m)
	if err != nil {
		return err
	}
	for k := range m {
		r.ips.Store(k, struct{}{})
	}
	return nil
}
