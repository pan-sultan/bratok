package config

import (
	"bufio"
	"fmt"
	"gontlm"
	"io"
	"os"
	"strings"
)

func ReadFromFile(filename string) (gontlm.Config, error) {
	cfg := gontlm.Config{}
	file, err := os.Open(filename)

	if err != nil {
		return cfg, err
	}

	defer file.Close()

	return Read(file)
}

func Read(reader io.Reader) (gontlm.Config, error) {
	cfg := gontlm.Config{}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}

		if err := setValue(&cfg, line); err != nil {
			return cfg, err
		}

		fmt.Println(line)
	}

	return cfg, nil
}
