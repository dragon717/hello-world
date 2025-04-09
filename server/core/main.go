package main

import (
	"Test/data_config"
	"fmt"
	"math/rand"
	"time"
)

var GParamCfg *data_config.XmlParameterConfig
var ActionParamCfg *data_config.XmlActionParameterConfig

var WorldMap *World

func main() {
	initmodel()
	GParamCfg = data_config.GetXmlParameterConfig()
	GParamCfg.LoadConfig()

	ActionParamCfg = data_config.GetXmlActionParameterConfig()
	ActionParamCfg.LoadConfig()

	WorldMap = NewWorld()

	npc1 := NewUser("小明", int32(rand.Intn(10000)), 18)
	tree1 := NewEntityTree("树1", int32(rand.Intn(10000)), 99)

	WorldMap.Gmap.SetLocation(1, 1, npc1)
	WorldMap.Gmap.SetLocation(1, 2, tree1)

	WorldMap.Gmap.Show()
	for {
		res := sendmsg(npc1, &ActionMsg{})
		fmt.Println(res)
		//WorldMap.Gmap.Show()
		time.Sleep(3 * time.Second)
	}

}
