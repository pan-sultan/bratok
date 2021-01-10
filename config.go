package gontlm

type Config struct {
	Username string
	PassNT   string
	Auth     string
	Domain   string
	Proxy    string
	NoProxy  []string
	Listen   []string
	Allow    []string
	Deny     []string
}
