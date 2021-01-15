package config

import (
	"fmt"
	"myproxy"
)

func Validate(cfg myproxy.Config) error {
	for _, p := range cfg.Proxies {
		switch p.Type {
		case myproxy.NTLM:
			if err := validateNTLM(p); err != nil {
				return err
			}
		case myproxy.Direct:
			if err := validateDirect(p); err != nil {
				return err
			}
		case myproxy.Socks5:
			if err := validateSocks5(p); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown proxy type: %s", p.Type)
		}
	}

	return nil
}

func validateNTLM(proxy myproxy.Proxy) error {
	return nil
}

func validateDirect(proxy myproxy.Proxy) error {
	return nil
}

func validateSocks5(proxy myproxy.Proxy) error {
	return nil
}
