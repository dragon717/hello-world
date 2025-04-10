package main

import (
	"flag"
	"fmt"
	"log"
)

var showMode = flag.Uint("mode", 0, "Show mode (0: Map, 1: Entities")

func main() {
	flag.Parse()
	serverAddr := flag.String("server", "localhost:8088", "Server address to connect to")
	flag.Parse()

	mon, err := NewMonitor(*serverAddr)
	if err != nil {
		log.Fatalf("Failed to create monitor: %v", err)
	}

	fmt.Println("Starting monitor...")
	mon.Start()
}

func ShowMode() uint {
	return *showMode
}
