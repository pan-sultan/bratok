package bratok

const (
	Direct = "direct"
	NTLM   = "ntlm"
	Socks  = "socks"
	HTTP   = "http"
)

type Config struct {
	Listen   []string
	Proxy    PAC
	Firewall Allower
}

// Proxy Auto-Configuration (PAC) like in mozilla projects
type PAC interface {
	Find(url string, host string) *Proxy
}

type Allower interface {
	Allow(ip string) bool
}

type Proxy struct {
	Type     string
	URL      string
	Name     string
	Username string
	Password string
	Domain   string
}
