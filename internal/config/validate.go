package config

import (
	"fmt"
	"bratok"
)

func Validate(cfg bratok.Config) error {
	for _, p := range cfg.Proxies {
		switch p.Type {
		case bratok.NTLM:
			if err := validateNTLM(p); err != nil {
				return err
			}
		case bratok.Direct:
			if err := validateDirect(p); err != nil {
				return err
			}
		case bratok.Socks5:
			if err := validateSocks5(p); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown proxy type: %s", p.Type)
		}
	}

	return nil
}

func validateNTLM(proxy bratok.Proxy) error {
	return nil
}

func validateDirect(proxy bratok.Proxy) error {
	return nil
}

func validateSocks5(proxy bratok.Proxy) error {
	return nil
}
