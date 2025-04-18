package main

import (
	"flag"
	"math/rand"
)

var WorldMap *World
var monitorServer *MonitorServer
var devMode = flag.Bool("dev", false, "Enable development mode (disable printing)")

func main() {
	flag.Parse() // 必须调用以解析flag参数
	InitCfg()
	initModelPool(16)

	WorldMap = NewWorld()
	WorldMap.AddEntity(NewUser("小明", int32(rand.Intn(10000)), 18, 1, 1))
	WorldMap.AddEntity(NewUser("小红", int32(rand.Intn(10000)), 18, 4, 1))
	WorldMap.AddEntity(NewUser("小蓝", int32(rand.Intn(10000)), 18, 1, 4))
	WorldMap.AddEntity(NewEntityTree("树1", int32(rand.Intn(10000)), 99, 1, 2))
	WorldMap.AddEntity(NewEntityTree("浆果灌木", int32(rand.Intn(10000)), 1, 7, 7))

	WorldMap.Gmap.Show()
	InitMonitorServer(":8088")
	select {}
}
