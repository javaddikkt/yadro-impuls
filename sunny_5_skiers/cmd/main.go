package main

import (
	"fmt"
	"log"
	"os"

	"sunny_5_skiers/internal/config"
	"sunny_5_skiers/internal/parser"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage: biathlon <config.json> <events_file> <output_file>")
		os.Exit(1)
	}

	cfgFile, _ := os.Open(os.Args[1])
	defer cfgFile.Close()

	var cfg config.Config
	if err := cfg.Decode(cfgFile); err != nil {
		log.Fatal(err)
	}

	if err := parser.ParseEvents(os.Args[2], os.Args[3], &cfg); err != nil {
		log.Fatal(err)
	}
}
