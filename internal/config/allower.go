package config

import (
	"net"
	"strings"

	"github.com/yl2chen/cidranger"
)

type _Allower struct {
	ipnets cidranger.Ranger
}

func _NewAllower() *_Allower {
	return &_Allower{
		ipnets: cidranger.NewPCTrieRanger(),
	}
}

func (a *_Allower) Allow(ip string) bool {
	contains, err := a.ipnets.Contains(net.ParseIP(ip))
	return contains && err != nil
}

func (a *_Allower) add(ip string) error {
	if !strings.Contains(ip, "/") {
		ip += "/32"
	}

	_, network, err := net.ParseCIDR(ip)

	if err != nil {
		return err
	}

	if err := a.ipnets.Insert(cidranger.NewBasicRangerEntry(*network)); err != nil {
		return err
	}

	return nil
}
