package main

import (
	"math/rand"
)

var WorldMap *World

func main() {
	InitCfg()
	initModelPool(16)

	WorldMap = NewWorld()
	WorldMap.AddEntity(NewUser("小明", int32(rand.Intn(10000)), 18, 1, 1))
	//WorldMap.AddEntity(NewUser("小红", int32(rand.Intn(10000)), 18, 4, 1))
	//WorldMap.AddEntity(NewUser("小蓝", int32(rand.Intn(10000)), 18, 1, 4))
	WorldMap.AddEntity(NewEntityTree("树1", int32(rand.Intn(10000)), 99, 1, 2))

	WorldMap.Gmap.Show()
	select {}

}
