package main

import (
	"fmt"
	"gontlm/internal/config"
	"log"
	"os"
)

func main() {
	cfg, err := config.ReadFromFile(os.Args[1])

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(cfg)
}
