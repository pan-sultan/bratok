package config

import (
	"bytes"
	"io"
	"myproxy"
	"os"

	"gopkg.in/yaml.v2"
)

func ReadFromFile(filename string) (myproxy.Config, error) {
	cfg := myproxy.Config{}
	file, err := os.Open(filename)

	if err != nil {
		return cfg, err
	}

	defer file.Close()

	return Read(file)
}

func Read(reader io.Reader) (myproxy.Config, error) {
	cfg := myproxy.Config{}
	data, err := readBytes(reader)

	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func readBytes(reader io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
