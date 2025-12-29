package ipblock

import (
	"encoding/json"
	"io"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

// Rules 表示被阻止ip的规则集合
type Rules struct {
	ips    sync.Map
	path   string
	m      sync.Mutex
	report Report
}

func (r *Rules) Add(ip, reason string) {
	r.ips.Store(ip, entry{Time: time.Now(), Reason: reason})
	if err := r.update(); err != nil {
		panic(err)
	}
	if r.report != nil {
		r.report.Report(ip, reason)
	}
}

type entry struct {
	Ip     string
	Time   time.Time
	Reason string
}

func (r *Rules) IsBlock(ip string) bool {
	_, ok := r.ips.Load(ip)
	return ok
}

// update 将新数据写入磁盘
func (r *Rules) update() error {
	// 加锁确保不会发生
	// g1和g2并发执行
	// g1是旧规则
	// g2是新规则并写入磁盘
	// g1写入旧规则到磁盘
	r.m.Lock()
	defer r.m.Unlock()
	fd, err := os.Create(r.path)
	if err != nil {
		return err
	}
	defer fd.Close()
	d := json.NewEncoder(fd)
	d.SetIndent("", "\t")

	// 获取最新规则
	var m = make([]entry, 0)
	r.ips.Range(func(key, value any) bool {
		e := value.(entry)
		m = append(m, entry{Ip: key.(string), Time: e.Time, Reason: e.Reason})
		return true
	})

	// 排序
	// 没有时间戳的按ip字典顺序排最前
	// 有时间戳的按时间先后排序
	slices.SortFunc(m, func(a, b entry) int {
		if !a.Time.IsZero() && !b.Time.IsZero() {
			return a.Time.Compare(b.Time)
		}
		if a.Time.IsZero() && b.Time.IsZero() {
			return strings.Compare(a.Ip, b.Ip)
		}
		if a.Time.IsZero() {
			return -1
		}
		if b.Time.IsZero() {
			return 1
		}
		panic("不可达路径")
	})

	err = d.Encode(&m)
	return err
}

// Init 读取磁盘中的旧数据
// 没有则创建新文件
func (r *Rules) Init(path string, report Report) error {
	r.path = path
	r.report = report
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	//尝试按旧方式读取规则
	err = r.try_old_decoer(fd)
	if err == nil {
		return nil
	}
	// 读取新方式保存的规则
	_, err = fd.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return r.try_new_decoer(fd)
}

func (r *Rules) try_old_decoer(fd *os.File) error {
	d := json.NewDecoder(fd)
	var m = make(map[string]struct{})
	err := d.Decode(&m)
	if err != nil {
		return err
	}
	for k := range m {
		r.ips.Store(k, entry{})
	}
	return nil
}

func (r *Rules) try_new_decoer(fd *os.File) error {
	d := json.NewDecoder(fd)
	var m = make([]entry, 0)
	err := d.Decode(&m)
	if err != nil {
		return err
	}
	for _, v := range m {
		r.ips.Store(v.Ip, v)
	}
	return nil
}
