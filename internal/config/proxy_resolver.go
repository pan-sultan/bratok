package config

import (
	"bratok"
)

func (r *_ProxyResolver) Find(url string, host string) *bratok.Proxy {
	for _, p := range r.proxies {
		if p.anyHosts || p.usefor.find(url, host) {
			return p.proxy
		}

		if p.exclude.length() != 0 {
			if !p.exclude.find(url, host) {
				return p.proxy
			}
		}

	}

	return nil
}

type _Proxy struct {
	name     string
	proxy    *bratok.Proxy
	usefor   *_PRStorage
	exclude  *_PRStorage
	anyHosts bool
}

type _ProxyResolver struct {
	proxies []*_Proxy
}

func _NewProxyResolver() *_ProxyResolver {
	return &_ProxyResolver{
		proxies: make([]*_Proxy, 0),
	}
}

func _NewProxy(proxy *bratok.Proxy) *_Proxy {
	return &_Proxy{
		name:    proxy.Name,
		proxy:   proxy,
		usefor:  _NewPRStorage(),
		exclude: _NewPRStorage(),
	}
}

func (r *_ProxyResolver) _Append(proxy *_Proxy) {
	r._PanicIfProxyExists(proxy.name)
	r.proxies = append(r.proxies, proxy)
}

func (r *_ProxyResolver) _PanicIfProxyExists(name string) {
	for _, p := range r.proxies {
		if name == p.name {
			panic("bug: proxy already added")
		}
	}
}
