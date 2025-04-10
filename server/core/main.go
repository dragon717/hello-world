package main

import (
	"flag"
	"math/rand"
	"time"
)

var WorldMap *World
var monitorServer *MonitorServer

func main() {
	InitCfg()
	flag.Parse() // 必须调用以解析参数

	initModelPool(16)

	WorldMap = NewWorld()
	WorldMap.AddEntity(NewUser("小明", int32(rand.Intn(10000)), 18, 1, 1))
	//WorldMap.AddEntity(NewUser("小红", int32(rand.Intn(10000)), 18, 4, 1))
	//WorldMap.AddEntity(NewUser("小蓝", int32(rand.Intn(10000)), 18, 1, 4))
	WorldMap.AddEntity(NewEntityTree("树1", int32(rand.Intn(10000)), 99, 1, 2))

	// Start monitor server
	var err error
	monitorServer, err = NewMonitorServer(":8088")
	if err != nil {
		panic(err)
	}
	monitorServer.Start()

	// Start update loop
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			monitorServer.Broadcast(monitorServer.GetMapData())
		}
	}()

	WorldMap.Gmap.Show()
	select {}
}
