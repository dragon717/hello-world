package main

import (
	"time"
)

type EntityTree struct {
	*Entity
}

// new one user
func NewEntityTree(name string, treeId int32, age uint32, x, y int) *EntityTree {
	e := &EntityTree{
		Entity: NewEntity(age, name, treeId, int32(GParamCfg.GetEntityTypeTree()), int32(age), x, y),
	}
	WorldMap.GEntityList[uint32(treeId)] = e
	WorldMap.GEntityTypeList[uint32(treeId)] = GParamCfg.GetEntityTypeTree()

	e.RegisterAction()
	go e.LifeProcess()
	return e
}

func (u *EntityTree) RegisterAction() {
	u.Register(ActionParamCfg.GetActionTypeCuttingDownTrees(), BeActionCuttingDownTrees)
}

func (u *EntityTree) LifeProcess() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			if u.HP <= 0 {
				WorldMap.Gmap.DeadChan <- u
				return
			}
			u.Age++
			u.HP++
			u.AddActionLog(&ActionLog{
				ActionType: ActionParamCfg.GetActionGrow(),
				Action:     "生长",
				Time:       WorldMap.Gmap.GlobalTime.GetTime(),
				Result:     "Age+1,HP+1",
			})
		}
	}
}
