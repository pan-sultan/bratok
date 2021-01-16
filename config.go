package bratok

import (
	"net"
	"unicode"

	"github.com/yl2chen/cidranger"
)

const (
	Direct = "direct"
	NTLM   = "ntlm"
	Socks  = "socks"
)

type Config struct {
	Listen        []string
	Allow         Hosts
	AutoconfigURL string
	Proxies       []*Proxy
	Rules         Hosts
}

type Proxy struct {
	Type string
	URL  string

	NTLM ProxyNTLM
}

type ProxyNTLM struct {
	PassNT   string
	Username string
	Domain   string
}

type Hosts struct {
	Allow    cidranger.Ranger
	Hostname map[string]*Proxy
}

type ProxyNetwork struct {
	Proxy *Proxy
	IPNet net.IPNet
}

func (p *ProxyNetwork) Network() net.IPNet {
	return p.IPNet
}

// host or ip
func (h *Hosts) FindProxy(host string) *Proxy {
	if unicode.IsDigit(rune(host[0])) { // host is IP
		containingNetworks, err := h.Allow.ContainingNetworks(net.ParseIP(host))

		if err != nil {
			panic(err)
		}

		if len(containingNetworks) > 0 {
			return containingNetworks[0].(*ProxyNetwork).Proxy
		}

		return nil
	}

	// host is hostname
	return h.Hostname[host]
}
