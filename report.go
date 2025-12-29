package ipblock

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

var _ Report = (*AbuseIPDB_Report)(nil)

// Report 向远程数据库报告1个ip
type Report interface {
	Report(ip, reason string)
}

// AbuseIPDB_Report 向AbuseIPDB报告1个ip
type AbuseIPDB_Report struct {
	Key string
}

func (r *AbuseIPDB_Report) Report(ip, reason string) {
	categories := 21
	if reason == "tls scan" {
		categories = 14
	}
	params := url.Values{}
	params.Add("ip", ip)
	params.Add("categories", fmt.Sprintf("%d", categories))
	params.Add("comment", reason)
	params.Add("timestamp", time.Now().Format(time.RFC3339))
	req, err := http.NewRequest(http.MethodPost, "https://api.abuseipdb.com/api/v2/report?"+params.Encode(), nil)
	if err != nil {
		slog.Error("report ip fail", "err", err)
		return
	}
	req.Header = make(http.Header)
	req.Header.Add("Key", r.Key)
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("report ip fail", "err", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		slog.Error("report ip fail", "status", resp.StatusCode)
		return
	}
	io.Copy(io.Discard, resp.Body)
}
