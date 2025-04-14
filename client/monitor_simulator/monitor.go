package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sort"
	"syscall"
)

type Monitor struct {
	conn     net.Conn
	mapData  [][]string
	mapCache [][]string
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
		mapCache: make([][]string, 0),
		entities: make(map[string]map[string]string),
	}, nil
}

func (m *Monitor) Start() {
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		m.receiveData()
		m.displayMap()
		m.displayEntities()
	}()

	<-sigChan
	m.conn.Close()
}

func (m *Monitor) receiveData() {
	buf := make([]byte, 100000)
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
		// Create new map data
		newMapData := make([][]string, len(mapData))
		for i, row := range mapData {
			if rowData, ok := row.([]interface{}); ok {
				newMapData[i] = make([]string, len(rowData))
				for j, cell := range rowData {
					if cellData, ok := cell.(string); ok {
						newMapData[i][j] = cellData
					}
				}
			}
		}

		// Check if map has changed
		hasChanged := false
		if len(m.mapCache) != len(newMapData) {
			hasChanged = true
		} else {
			for i := range newMapData {
				if len(m.mapCache[i]) != len(newMapData[i]) {
					hasChanged = true
					break
				}
				for j := range newMapData[i] {
					if m.mapCache[i][j] != newMapData[i][j] {
						hasChanged = true
						break
					}
				}
				if hasChanged {
					break
				}
			}
		}

		if hasChanged {
			m.mapData = newMapData
			m.mapCache = make([][]string, len(newMapData))
			for i := range newMapData {
				m.mapCache[i] = make([]string, len(newMapData[i]))
				copy(m.mapCache[i], newMapData[i])
			}
			// Only display when map data changes
			m.displayMap()
		}
	}

	// Process entity data
	if entities, ok := data["entities"].(map[string]interface{}); ok {
		m.entities = make(map[string]map[string]string)
		for id, entity := range entities {
			if entityData, ok := entity.(map[string]interface{}); ok {
				m.entities[id] = make(map[string]string)
				for k, v := range entityData {
					m.entities[id][k] = fmt.Sprintf("%+v", v)
				}
			}
		}
	}
	m.displayEntities()
}

func (m *Monitor) displayMap() {
	if ShowMode() != 0 {
		return
	}

	// Only display map if it has changed
	if len(m.mapData) > 0 {
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
	}
}

func (m *Monitor) displayEntities() {
	if ShowMode() != 1 {
		return
	}

	// Display entities
	fmt.Println("\n=== Entities ===")

	// Sort the keys
	keys := make([]string, 0, len(m.entities))
	for id := range m.entities {
		keys = append(keys, id)
	}
	sort.Strings(keys)

	for _, id := range keys {
		entity := m.entities[id]
		fmt.Printf("ID: %s\n", id)

		// Sort the entity keys
		entityKeys := make([]string, 0, len(entity))
		for k := range entity {
			entityKeys = append(entityKeys, k)
		}
		sort.Strings(entityKeys)

		for _, k := range entityKeys {
			v := entity[k]
			fmt.Printf("  %s: %s\n", k, v)
		}
		fmt.Println()
	}
}
