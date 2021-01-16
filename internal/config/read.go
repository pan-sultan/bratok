package config

import (
	"bratok"
	"bytes"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

func ReadFromFile(filename string) (*bratok.Config, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return Read(file)
}

func Read(reader io.Reader) (*bratok.Config, error) {
	data, err := readBytes(reader)

	if err != nil {
		return nil, err
	}

	ycfg := yamlCfg{}
	if err := yaml.Unmarshal(data, &ycfg); err != nil {
		return nil, err
	}

	validate(ycfg)

	return yaml2Config(ycfg), nil
}

func readBytes(reader io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
