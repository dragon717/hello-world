package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	serverAddr := flag.String("server", "localhost:8088", "Server address to connect to")
	flag.Parse()

	mon, err := NewMonitor(*serverAddr)
	if err != nil {
		log.Fatalf("Failed to create monitor: %v", err)
	}

	fmt.Println("Starting monitor...")
	mon.Start()
}
