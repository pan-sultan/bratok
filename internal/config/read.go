package config

import (
	"bratok"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

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

	ycfg := yConfig{}
	if err := yaml.UnmarshalStrict(data, &ycfg); err != nil {
		return nil, convertYamlErr(err)
	}

	if err := validate(ycfg); err != nil {
		return nil, err
	}

	return yaml2Config(ycfg)
}

func readBytes(reader io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func convertYamlErr(err error) error {
	s := err.Error()
	if !(strings.Contains(s, "not found") && strings.Contains(s, "field")) {
		return err
	}

	ss := strings.Split(s, ":")

	for _, s := range ss {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "line ") {
			return fmt.Errorf("%s: unknown parameter", s)
		}
	}

	return err
}
