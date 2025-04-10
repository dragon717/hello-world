package main

import (
	"encoding/json"
	"fmt"
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
		}
	}
}

// 获取周围信息
func (m *Wmap) GetAroundInfo(x, y int, size int) string {
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
				ejson, _ := json.Marshal(entityInterface.GetInfo())
				infoMap[entityInterface.GetName()] = string(ejson)
			}
		}
	}
	info := fmt.Sprintf("%+v", infoMap)
	return string(info)
}
