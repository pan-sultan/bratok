package config

import "bratok"

type yamlCfg struct {
	Listen        []string            `yaml:"listen"`
	Allow         []string            `yaml:"allow"`
	AutoconfigURL string              `yaml:"autoconfig_url"`
	Proxies       []*bratok.Proxy     `yaml:"proxies"`
	Rules         map[string][]string `yaml:"rules"`
}

type Proxy struct {
	Type string `yaml:"type"`
	URL  string `yaml:"url"`

	NTLM ProxyNTLM `yaml:",inline"`
}

type ProxyNTLM struct {
	PassNT   string `yaml:"passhash"`
	Username string `yaml:"user"`
	Domain   string `yaml:"domain"`
}

func yaml2Config(y yamlCfg) *bratok.Config {
	cfg := new(bratok.Config)

	cfg.Listen = y.Listen
	//cfg.Allow = y.Allow
	cfg.AutoconfigURL = y.AutoconfigURL
	cfg.Proxies = y.Proxies
	return cfg
}
