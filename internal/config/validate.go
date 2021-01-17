package config

import (
	"bratok"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

func validate(y yConfig) error {
	if err := validateListen(y.Listen); err != nil {
		return fmt.Errorf("%s: %v", yConfigFieldTagName("Listen"), err)
	}

	if err := validateAllow(y.Allow); err != nil {
		return fmt.Errorf("%s: %v", yConfigFieldTagName("Allow"), err)
	}

	if err := validateProxies(y.Proxies); err != nil {
		return fmt.Errorf("%s: %v", yConfigFieldTagName("Proxies"), err)
	}

	if err := validateRules(y); err != nil {
		return fmt.Errorf("%s: %v", yConfigFieldTagName("Rules"), err)
	}

	return nil
}

func validateListen(listen []string) error {
	if len(listen) == 0 {
		return errors.New("cannot be empty")
	}

	for _, listenIP := range listen {
		ipPort := strings.Split(listenIP, ":")

		if len(ipPort) != 2 {
			return fmt.Errorf("invalid listen format '%s', it must be <IPv4>:<port>", listenIP)
		}

		ip, port := ipPort[0], ipPort[1]

		if err := validateIPv4(ip); err != nil {
			return err
		}

		if err := validatePort(port); err != nil {
			return err
		}
	}

	return nil
}

func validatePort(port string) error {
	if n, err := strconv.Atoi(port); err != nil || n < 1 || n > 65535 {
		return fmt.Errorf("invalid port: '%s'", port)
	}

	return nil
}

func validateAllow(allow []string) error {
	if len(allow) == 0 {
		return errors.New("cannot be empty")
	}

	for _, ip := range allow {
		if err := validateIPv4(ip); err != nil {
			return err
		}
	}

	return nil
}

func validateIPv4(host string) error {
	if net.ParseIP(host) != nil && !strings.Contains(host, ":") {
		return nil
	}

	return errors.New("invalid IPv4 address")
}

func validateHostname(host string) error {
	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	if re.MatchString(host) {
		return nil
	}

	return fmt.Errorf("invalid hostname: %s", host)
}

func validateProxies(proxies []*yProxy) error {
	if len(proxies) == 0 {
		return configFieldEmptyErr("Proxies")
	}

	for _, p := range proxies {
		if err := validateProxyDependsOnType(*p); err != nil {
			return err
		}
	}

	return nil
}

func validateRules(y yConfig) error {
	proxies := y.yProxyMap()

	if err := validateIfAllRulesSpecified(proxies, y.Rules); err != nil {
		return err
	}

	if err := validateRulesPatterns(y.Rules); err != nil {
		return err
	}

	return nil
}

func validateRulesPatterns(rules map[string][]string) error {
	for _, patterns := range rules {
		for _, p := range patterns {
			if p == "*" {
				continue
			}

			err := func() error {
				p = strings.TrimPrefix(p, "!")

				if isIP(p) {
					if strings.Contains(p, "/") {
						if _, _, err := net.ParseCIDR(p); err != nil {
							return err
						}
					} else if err := validateIPv4(p); err != nil {
						return err
					}
				} else {
					if err := validateHostname(p); err != nil {
						return err
					}
				}

				return nil
			}()

			if err != nil {
				return fmt.Errorf("wrong pattern: %s: %v", p, err)
			}
		}
	}

	return nil
}

func validateIfAllRulesSpecified(proxies map[string]*yProxy, rules map[string][]string) error {
	err := fmt.Errorf(
		"%s must be specified for each %s",
		yConfigFieldTagName("Rules"),
		yConfigFieldTagName("Proxies"),
	)

	if len(rules) != len(proxies) {
		return err
	}

	for name := range proxies {
		if _, found := rules[name]; !found {
			return err
		}
	}

	return nil
}

func validateProxyDependsOnType(p yProxy) error {
	err := func() error {
		if err := validateProxy(p); err != nil {
			return err
		}

		switch p.Type {
		case bratok.Direct:
			if err := validateDirect(p); err != nil {
				return err
			}
		case bratok.NTLM:
			if err := validateNTLM(p); err != nil {
				return err
			}
		case bratok.HTTP:
			if err := validateHTTP(p); err != nil {
				return err
			}
		case bratok.Socks:
			if err := validateSocks5(p); err != nil {
				return err
			}
		default:
			return errors.New("unknown proxy type")
		}

		return nil
	}()

	if err != nil {
		return fmt.Errorf("%s: %v", p.Type, err)
	}

	return nil
}

func validateProxy(p yProxy) error {
	if p.Name == "" {
		return errors.New("name empty")
	}

	if p.Type != bratok.Direct {
		if err := validateProxyURL(p.URL); err != nil {
			return err
		}
	}

	return nil
}

func validateProxyURL(url string) error {
	if url == "" {
		return proxyFieldEmptyErr("URL")
	}

	s := strings.Split(url, ":")
	if len(s) != 2 {
		return fmt.Errorf("URL must be <hostname or IP>:<port>: %s", url)
	}

	if isIP(s[0]) {
		if err := validateIPv4(s[0]); err != nil {
			return err
		}
	} else {
		if err := validateHostname(s[0]); err != nil {
			return err
		}
	}

	if err := validatePort(s[1]); err != nil {
		return err
	}

	return nil
}

func validateHTTP(p yProxy) error {
	if p.Password == "" {
		return proxyFieldEmptyErr("Password")
	}

	if p.Username == "" {
		return proxyFieldEmptyErr("Username")
	}

	p.Type = ""
	p.Name = ""
	p.Password = ""
	p.URL = ""
	p.Username = ""

	if err := validateForUselessFields(p); err != nil {
		return err
	}

	return nil
}

func validateNTLM(p yProxy) error {
	if p.Domain == "" {
		return proxyFieldEmptyErr("Domain")
	}

	if p.PassNT == "" {
		return proxyFieldEmptyErr("PassNT")
	}

	if p.Username == "" {
		return proxyFieldEmptyErr("Username")
	}

	p.Type = ""
	p.Domain = ""
	p.Name = ""
	p.PassNT = ""
	p.URL = ""
	p.Username = ""

	if err := validateForUselessFields(p); err != nil {
		return err
	}

	return nil
}

func validateDirect(p yProxy) error {
	p.Type = ""
	p.Name = ""

	if err := validateForUselessFields(p); err != nil {
		return err
	}

	return nil
}

func validateSocks5(p yProxy) error {
	p.Type = ""
	p.Name = ""
	p.URL = ""

	if err := validateForUselessFields(p); err != nil {
		return err
	}

	return nil
}

func proxyFieldEmptyErr(fieldname string) error {
	return errors.New(yProxyFieldTagName(fieldname) + " empty")
}

func configFieldEmptyErr(fieldname string) error {
	return errors.New(yConfigFieldTagName(fieldname) + " empty")
}

func validateForUselessFields(p yProxy) error {
	if data, err := yaml.Marshal(p); strings.TrimSpace(string(data)) != "{}" || err != nil {
		if err != nil {
			panic("bug")
		}

		return fmt.Errorf("useless options specified: %s", string(data))
	}

	return nil
}
