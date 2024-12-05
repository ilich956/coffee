package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// CHANGE LOGGING
var (
	Port *string
	Dir  *string
)

func ParseFlags() {
	Port = flag.String("port", "8080", "Port number")
	Dir = flag.String("dir", "data", "Path to the data directory")
	help := flag.Bool("help", false, "Show help screen")
	flag.Parse()

	if *help {
		fmt.Println(`Coffee Shop Management System

		Usage:
		  hot-coffee [--port <N>] [--dir <S>] 
		  hot-coffee --help
		
		Options:
		  --help       Show this screen.
		  --port N     Port number.
		  --dir S      Path to the data directory.`)

		os.Exit(0)
	}

	if err := validateDir(); err != nil {
		log.Fatal(err)
	}

	if err := validatePort(); err != nil {
		log.Fatal(err)
	}
}

func validatePort() error {
	port, err := strconv.Atoi(*Port)
	if err != nil {
		return fmt.Errorf("port should be number")
	}

	if port < 1024 || port > 49151 {
		return fmt.Errorf("invalid port, must be between 1024 and 49151")
	}

	return nil
}

// CHANGE VALIDATION
func validateDir() error {
	// if *Dir == "" || *Dir == "pkg" {
	// 	return fmt.Errorf("forbidden dir")
	// }
	if strings.Contains(*Dir, "cmd") || strings.Contains(*Dir, "internal") || strings.Contains(*Dir, "models") {
		return fmt.Errorf("forbidden dir")
	}
	return nil
}
