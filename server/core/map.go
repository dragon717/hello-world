package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"
)

type Wmap struct {
	Map        [][]*Block
	Size       uint64
	DeadChan   chan EntityInterface
	GlobalTime *MapTime
}

func NewMap(size uint64) *Wmap {
	mp := &Wmap{
		Map:      make([][]*Block, size),
		Size:     size,
		DeadChan: make(chan EntityInterface),
		GlobalTime: &MapTime{
			Month: time.January,
			Day:   1,
			Hour:  8,
		},
	}
	for x, _ := range mp.Map {
		mp.Map[x] = make([]*Block, size)
		for y, _ := range mp.Map[x] {
			mp.Map[x][y] = NewBlock()
		}
	}
	go mp.TimeProcess()
	return mp
}

func (m *Wmap) SetLocation(x, y int, e EntityInterface) {
	if m.Map[x][y] == nil {
		m.Map[x][y] = &Block{
			EntityList: make([]EntityInterface, 0),
		}
	}
	e.SetX(uint32(x))
	e.SetY(uint32(y))
	m.Map[x][y].AddEntity(e)
}
func (m *Wmap) MoveLocation(oldX, oldY, x, y int, e EntityInterface) {
	for i, _ := range m.Map[oldX][oldY].EntityList {
		if m.Map[oldX][oldY].EntityList[i].GetId() == e.GetId() {
			m.Map[oldX][oldY].EntityList = append(m.Map[oldX][oldY].EntityList[:i], m.Map[oldX][oldY].EntityList[i+1:]...)
			break
		}
	}
	if m.Map[x][y] == nil {
		m.Map[x][y] = &Block{
			EntityList: make([]EntityInterface, 0),
		}
	}
	e.SetX(uint32(x))
	e.SetY(uint32(y))
	m.Map[x][y].AddEntity(e)
}

func (w *Wmap) Show() {
	for _, blocks := range w.Map {
		for _, block := range blocks {
			for _, entityInterface := range block.EntityList {
				fmt.Print(entityInterface.GetId(), ":", entityInterface.GetName(), ":", entityInterface.GetType())
			}
			fmt.Print("|")
		}
		fmt.Println()
	}
}

func (m *Wmap) TimeProcess() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case entity := <-m.DeadChan:
			entity.SetHP(0)
			entity.AddActionLog(&ActionLog{
				ActionType: 0,
				Action:     "死亡",
				Time:       WorldMap.Gmap.GlobalTime.GetTime(),
			})
			delete(WorldMap.GEntityList, uint32(entity.GetId()))
			delete(WorldMap.GEntityTypeList, uint32(entity.GetId()))
			fmt.Println(entity, " dead")
		case <-ticker.C:
			m.GlobalTime.AddHour()

			m.GenerateRandItem()
		}
	}
}

// 获取周围信息
func (m *Wmap) GetAroundInfo(x, y int, size int, selfId int32) string {
	var infoMap map[string]string
	infoMap = make(map[string]string)
	for i := x - size; i <= x+size; i++ {
		for j := y - size; j <= y+size; j++ {
			if i < 0 || j < 0 || i >= int(m.Size) || j >= int(m.Size) {
				continue
			}
			if m.Map[i][j] == nil {
				continue
			}
			for _, entityInterface := range m.Map[i][j].EntityList {
				if entityInterface.GetId() == selfId {
					continue
				}
				ejson, _ := json.Marshal(entityInterface.GetInfo(false))
				infoMap[entityInterface.GetName()] = string(ejson)
			}
			for _, item := range m.Map[i][j].ItemList {
				infoMap[item.Name] = fmt.Sprintf("%+v", item)
			}
		}
	}
	info := fmt.Sprintf("%+v", infoMap)
	return string(info)
}
func (m *Wmap) GenerateRandItem() {
	it := &Item{
		Num: 1,
	}
	ranx := rand.IntN(int(m.Size - 1))
	rany := rand.IntN(int(m.Size - 1))
	if m.Map[ranx][rany].ItemList == nil {
		m.Map[ranx][rany].ItemList = make([]*Item, 0)
	}
	ran := rand.IntN(3)
	switch ran {
	case 0:
		it.ID = ItemParamCfg.GetItemBranch()
		it.Name = ItemCfg.GetById(int(ItemParamCfg.GetItemBranch())).Name
		it.Type = ItemtypeParamCfg.GetItemTypeMaterial()
	case 1:
		it.ID = ItemParamCfg.GetItemPebble()
		it.Name = ItemCfg.GetById(int(ItemParamCfg.GetItemPebble())).Name
		it.Type = ItemtypeParamCfg.GetItemTypeMaterial()
	case 2:
		it.ID = ItemParamCfg.GetItemBerry()
		it.Name = ItemCfg.GetById(int(ItemParamCfg.GetItemBerry())).Name
		it.Type = ItemtypeParamCfg.GetItemTypeMaterial()
	}
	m.Map[ranx][rany].ItemList = append(m.Map[ranx][rany].ItemList, it)

}
