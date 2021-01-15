package myproxy

const (
	Direct = "direct"
	NTLM   = "ntlm"
	Socks5 = "socks5"
)

type Config struct {
	NoProxy []string
	Listen  []string `yaml:"listen"`
	Allow   []string `yaml:"allow"`
	Deny    []string `yaml:"deny"`
	Proxies []Proxy  `yaml:"proxies"`
}

type Proxy struct {
	Type   string   `yaml:"type"`
	URL    string   `yaml:"url"`
	UseFor []string `yaml:"usefor"`
	Ignore []string `yaml:"ignore"`

	ProxyNTLM `yaml:",inline"`
}

type ProxyNTLM struct {
	PassNT   string `yaml:"passhash"`
	Username string `yaml:"user"`
	Domain   string `yaml:"domain"`
}
