package main

import (
	"bratok/internal/config"
	"fmt"
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
