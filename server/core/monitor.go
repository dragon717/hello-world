package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type MonitorServer struct {
	listener net.Listener
	clients  []net.Conn
	mu       sync.Mutex
}

func NewMonitorServer(addr string) (*MonitorServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to start monitor server: %v", err)
	}

	return &MonitorServer{
		listener: listener,
		clients:  make([]net.Conn, 0),
	}, nil
}

func (ms *MonitorServer) Start() {
	go ms.acceptConnections()
}

func (ms *MonitorServer) acceptConnections() {
	for {
		conn, err := ms.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		ms.mu.Lock()
		ms.clients = append(ms.clients, conn)
		ms.mu.Unlock()

		fmt.Printf("New monitor client connected: %s\n", conn.RemoteAddr())
	}
}

func (ms *MonitorServer) Broadcast(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Remove disconnected clients
	var activeClients []net.Conn
	for _, client := range ms.clients {
		_, err := client.Write(jsonData)
		if err != nil {
			fmt.Printf("Error sending data to client: %v\n", err)
			client.Close()
			continue
		}
		activeClients = append(activeClients, client)
	}
	ms.clients = activeClients
}

func (ms *MonitorServer) GetMapData() interface{} {
	data := make(map[string]interface{})

	// Get map data
	mapData := make([][]string, len(WorldMap.Gmap.Map))
	for i, row := range WorldMap.Gmap.Map {
		mapData[i] = make([]string, len(row))
		for j, block := range row {
			if len(block.EntityList) > 0 {
				entity := block.EntityList[0]
				mapData[i][j] = fmt.Sprintf("%s(%d)", entity.GetName(), entity.GetId())
			} else {
				mapData[i][j] = ""
			}
		}
	}
	data["map"] = mapData

	// Get entity data
	entities := make(map[string]map[string]string)
	for _, entity := range WorldMap.GEntityList {
		entities[fmt.Sprintf("%d", entity.GetId())] = entity.GetInfo()
	}
	data["entities"] = entities

	return data
}
