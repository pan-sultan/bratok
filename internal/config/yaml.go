package config

import (
	"bratok"
	"fmt"
	"reflect"
	"strings"
)

type yConfig struct {
	Listen        []string            `yaml:"listen"`
	Allow         []string            `yaml:"allow"`
	AutoconfigURL string              `yaml:"autoconfig_url"`
	Proxies       []*yProxy           `yaml:"proxies"`
	Rules         map[string][]string `yaml:"rules"`
}

type yProxy struct {
	Type     string `yaml:"type,omitempty"`
	URL      string `yaml:"url,omitempty"`
	Name     string `yaml:"name,omitempty"`
	PassNT   string `yaml:"passhash,omitempty"`
	Password string `yaml:"password,omitempty"`
	Username string `yaml:"user,omitempty"`
	Domain   string `yaml:"domain,omitempty"`
}

func (y yConfig) yProxyMap() map[string]*yProxy {
	proxies := make(map[string]*yProxy)
	for _, p := range y.Proxies {
		proxies[p.Name] = p
	}

	return proxies
}

func yProxyFieldTagName(fieldName string) string {
	return structFieldTagName(&yProxy{}, fieldName)
}

func yConfigFieldTagName(fieldName string) string {
	return structFieldTagName(&yConfig{}, fieldName)
}

func structFieldTagName(s interface{}, fieldName string) string {
	f, ok := reflect.TypeOf(s).Elem().FieldByName(fieldName)
	if !ok {
		panic("bug: " + fieldName)
	}

	tag := string(f.Tag)
	a := strings.Split(tag, ",")
	return strings.Split(a[0], "\"")[1]
}

func yaml2Config(y yConfig) (*bratok.Config, error) {
	cfg := new(bratok.Config)
	cfg.Listen = y.Listen

	if err := fillAllower(cfg, y); err != nil {
		return nil, fmt.Errorf("%s: %v", yConfigFieldTagName("Allow"), err)
	}

	if err := fillProxies(cfg, y); err != nil {
		return nil, fmt.Errorf("%s: %v", yConfigFieldTagName("Proxies"), err)
	}

	return cfg, nil
}

func fillAllower(cfg *bratok.Config, y yConfig) error {
	allower := _NewFirewall()
	for _, ip := range y.Allow {
		if err := allower.add(ip); err != nil {
			return err
		}
	}

	cfg.Firewall = allower
	return nil
}

func fillProxies(cfg *bratok.Config, y yConfig) error {
	proxies := y.yProxyMap()
	pr := _NewProxyResolver()

	for name, patterns := range y.Rules {
		proxy, found := proxies[name]

		if !found || proxy == nil {
			return fmt.Errorf("not found: %s", name)
		}

		storage, err := createProxyStorage(proxy, patterns)
		if err != nil {
			return err
		}

		pr._Append(storage)
	}

	cfg.Proxy = pr
	return nil
}

func createProxyStorage(yproxy *yProxy, patterns []string) (*_Proxy, error) {
	storage := _NewProxy(_CreateProxy(*yproxy))

	for _, pattern := range patterns {
		if pattern == "*" {
			storage.all = true
			continue
		}

		exclude := strings.HasPrefix(pattern, "!")
		var s *_PRStorage

		if exclude {
			pattern = pattern[1:]
			s = storage.exclude
		} else {
			s = storage.usefor
		}

		if err := s.add(pattern); err != nil {
			return nil, err
		}
	}

	return storage, nil

}

func _CreateProxy(y yProxy) *bratok.Proxy {
	switch y.Type {
	case bratok.NTLM:
		return _NewNTLM(y)
	case bratok.Direct:
		return _NewDirect(y)
	case bratok.Socks:
		return _NewSocks(y)
	case bratok.HTTP:
		return _NewHTTP(y)
	}

	panic("bug")
}

func _NewNTLM(y yProxy) *bratok.Proxy {
	return &bratok.Proxy{
		Type:     bratok.NTLM,
		Name:     y.Name,
		URL:      y.URL,
		Domain:   y.Domain,
		Password: y.PassNT,
		Username: y.Username,
	}
}

func _NewDirect(y yProxy) *bratok.Proxy {
	return &bratok.Proxy{
		Type: bratok.Direct,
		Name: y.Name,
	}
}

func _NewSocks(y yProxy) *bratok.Proxy {
	return &bratok.Proxy{
		Type: bratok.NTLM,
		Name: y.Name,
		URL:  y.URL,
	}
}

func _NewHTTP(y yProxy) *bratok.Proxy {
	return &bratok.Proxy{
		Type:     bratok.HTTP,
		Name:     y.Name,
		URL:      y.URL,
		Password: y.Password,
		Username: y.Username,
	}
}
