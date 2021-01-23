package config

import (
	"net"
	"strings"

	"github.com/yl2chen/cidranger"
)

type _Firewall struct {
	ipnets cidranger.Ranger
}

func _NewFirewall() *_Firewall {
	return &_Firewall{
		ipnets: cidranger.NewPCTrieRanger(),
	}
}

func (a *_Firewall) Allow(ip string) bool {
	contains, err := a.ipnets.Contains(net.ParseIP(ip))
	return contains && err != nil
}

func (a *_Firewall) add(ip string) error {
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
