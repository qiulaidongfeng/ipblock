// Package 提供持久化封禁ip
package ipblock

import (
	"strings"
)

var match_all_path_rule = []string{
	"/env",
}

var match_prefix_path_rule = []string{
	"/webadmin",
	"/wp-content",
	"/admin",
	"/cgi-bin",
	"/config",
	"/.git",
	"/..",
}

// MayAttack 通过检查urlpath判断是否可能是恶意攻击
func MayAttack(urlpath string) bool {
	for _, s := range match_all_path_rule {
		if urlpath == s {
			return true
		}
	}
	for _, s := range match_prefix_path_rule {
		if strings.HasPrefix(urlpath, s) {
			return true
		}
	}
	return false
}
