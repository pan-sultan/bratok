package config

import "unicode"

func isIP(host string) bool {
	return unicode.IsDigit(rune(host[0]))
}
