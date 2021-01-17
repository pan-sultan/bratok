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
	Firewall Firewall
}

// Proxy Auto-Configuration (PAC) like in mozilla projects
type PAC interface {
	Find(url string, host string) *Proxy
}

type Firewall interface {
	Allow(ip string) bool
}

// Plain Proxy struct used for all proxy types
type Proxy struct {
	Type     string
	URL      string
	Name     string
	Username string
	Password string
	Domain   string
}
