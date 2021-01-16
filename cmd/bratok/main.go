package main

import (
	"fmt"
	"log"
	"bratok"
	"bratok/internal/config"
	"os"
)

func main() {
	cfg, err := readConfig(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(cfg)
}

func readConfig(filename string) (bratok.Config, error) {
	cfg, err := config.ReadFromFile(os.Args[1])

	if err != nil {
		return cfg, err
	}

	if err := config.Validate(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
