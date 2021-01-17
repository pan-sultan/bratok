package config

import (
	"net"
	"strings"

	"github.com/yl2chen/cidranger"
)

type _PRStorage struct {
	ipnets   cidranger.Ranger
	hostname map[string]struct{}
}

func _NewPRStorage() *_PRStorage {
	return &_PRStorage{
		ipnets:   cidranger.NewPCTrieRanger(),
		hostname: make(map[string]struct{}),
	}
}

func (r *_PRStorage) length() int {
	return len(r.hostname) + r.ipnets.Len()
}

func (r *_PRStorage) find(url string, host string) bool {
	if isIP(host) {
		nets, err := r.ipnets.ContainingNetworks(net.ParseIP(host))
		return err == nil && len(nets) > 0
	}

	_, found := r.hostname[host]
	return found
}

func (r *_PRStorage) add(host string) error {
	if isIP(host) {
		return r.addAsCIDR(host)
	}

	r.addAsHost(host)
	return nil
}

func (r *_PRStorage) addAsCIDR(ip string) error {
	if !strings.Contains(ip, "/") {
		ip += "/32"
	}

	_, network, err := net.ParseCIDR(ip)

	if err != nil {
		return err
	}

	if err := r.ipnets.Insert(cidranger.NewBasicRangerEntry(*network)); err != nil {
		return err
	}

	return nil
}

func (r *_PRStorage) addAsHost(host string) {
	r.hostname[host] = struct{}{}
}
