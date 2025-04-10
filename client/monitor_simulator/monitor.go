package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Monitor struct {
	conn     net.Conn
	mapData  [][]string
	entities map[string]map[string]string
}

func NewMonitor(serverAddr string) (*Monitor, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	return &Monitor{
		conn:     conn,
		mapData:  make([][]string, 0),
		entities: make(map[string]map[string]string),
	}, nil
}

func (m *Monitor) Start() {
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go m.receiveData()
	go m.displayLoop()

	<-sigChan
	m.conn.Close()
}

func (m *Monitor) receiveData() {
	buf := make([]byte, 4096)
	for {
		n, err := m.conn.Read(buf)
		if err != nil {
			fmt.Printf("Error reading from server: %v\n", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(buf[:n], &data); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			continue
		}

		m.processData(data)
	}
}

func (m *Monitor) processData(data map[string]interface{}) {
	// Process map data
	if mapData, ok := data["map"].([]interface{}); ok {
		m.mapData = make([][]string, len(mapData))
		for i, row := range mapData {
			if rowData, ok := row.([]interface{}); ok {
				m.mapData[i] = make([]string, len(rowData))
				for j, cell := range rowData {
					if cellData, ok := cell.(string); ok {
						m.mapData[i][j] = cellData
					}
				}
			}
		}
	}

	// Process entity data
	if entities, ok := data["entities"].(map[string]interface{}); ok {
		m.entities = make(map[string]map[string]string)
		for id, entity := range entities {
			if entityData, ok := entity.(map[string]interface{}); ok {
				m.entities[id] = make(map[string]string)
				for k, v := range entityData {
					m.entities[id][k] = fmt.Sprintf("%v", v)
				}
			}
		}
	}
}

func (m *Monitor) displayLoop() {
	for {
		// Clear screen
		fmt.Print("\033[H\033[2J")

		// Display map
		fmt.Println("=== Map ===")
		for _, row := range m.mapData {
			for _, cell := range row {
				if cell == "" {
					fmt.Print(".")
				} else {
					fmt.Print(cell)
				}
				fmt.Print(" ")
			}
			fmt.Println()
		}

		// Display entities
		fmt.Println("\n=== Entities ===")
		for id, entity := range m.entities {
			fmt.Printf("ID: %s\n", id)
			for k, v := range entity {
				fmt.Printf("  %s: %s\n", k, v)
			}
			fmt.Println()
		}

		// Wait a bit before refreshing
		time.Sleep(1 * time.Second)
	}
}
