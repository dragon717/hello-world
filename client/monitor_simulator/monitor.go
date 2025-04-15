package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

type Monitor struct {
	conn        net.Conn
	serverAddr  string
	mapData     [][]string
	mapCache    [][]string
	entities    map[string]map[string]string
	entityCache map[string]map[string]string
}

func NewMonitor(serverAddr string) (*Monitor, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	return &Monitor{
		conn:        conn,
		serverAddr:  serverAddr,
		mapData:     make([][]string, 0),
		mapCache:    make([][]string, 0),
		entities:    make(map[string]map[string]string),
		entityCache: make(map[string]map[string]string),
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

func (m *Monitor) reconnect() error {
	if m.conn != nil {
		m.conn.Close()
	}
	
	var err error
	m.conn, err = net.Dial("tcp", m.serverAddr)
	if err != nil {
		return fmt.Errorf("reconnect failed: %v", err)
	}
	fmt.Println("Successfully reconnected to server")
	return nil
}

func (m *Monitor) receiveData() {
	buf := make([]byte, 100000)
	for {
		if m.conn == nil {
			if err := m.reconnect(); err != nil {
				fmt.Printf("Reconnect error: %v, retrying in 3 seconds...\n", err)
				time.Sleep(3 * time.Second)
				continue
			}
		}

		n, err := m.conn.Read(buf)
		if err != nil {
			fmt.Printf("Connection error: %v, attempting to reconnect...\n", err)
			m.conn = nil
			time.Sleep(3 * time.Second)
			continue
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

	// Update entityCache
	m.entityCache = make(map[string]map[string]string)
	for id, entity := range m.entities {
		m.entityCache[id] = make(map[string]string)
		for k, v := range entity {
			m.entityCache[id][k] = v
		}
	}
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

	// Sort the keys
	keys := make([]string, 0, len(m.entities))
	for id := range m.entities {
		keys = append(keys, id)
	}
	sort.Strings(keys)

	for _, id := range keys {
		entity := m.entities[id]
		cachedEntity, ok := m.entityCache[id]
		if !ok {
			fmt.Println("\n=== Entities complete packege ===")
			fmt.Printf("ID: %s (New)\n", id)
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
			continue
		}

		// Compare entity properties
		changed := false
		entityKeys := make([]string, 0, len(entity))
		for k := range entity {
			entityKeys = append(entityKeys, k)
		}
		sort.Strings(entityKeys)
		for _, k := range entityKeys {
			v := entity[k]
			cachedV, ok := cachedEntity[k]
			if !ok || v != cachedV {
				if !changed {
					fmt.Printf("ID: %s (Updated)\n", id)
					changed = true
				}
				fmt.Println("\n=== Entities change ===")
				fmt.Printf("id: %s name: %s", entity["id"], entity["name"])
				cachedValue := cachedEntity[k]
				var lastValues []map[string]string
				json.Unmarshal([]byte(cachedValue), &lastValues)
				if len(lastValues) > 0 {
					var lastValue = lastValues[len(lastValues)-1]
					fmt.Printf("  %+v \n", lastValue)
				} else {
					fmt.Printf("  %s: %s (was: %s)\n", k, v, cachedValue)
				}
			}
		}
		if changed {
			fmt.Println()
		}
	}
}
