package main

import (
	"time"
)

type EntityUser struct {
	Tool  []string
	Money uint64
	*Entity
}

func NewUser(name string, npcId int32, age uint32) *EntityUser {
	e := &EntityUser{
		Tool:   make([]string, 0),
		Money:  0,
		Entity: NewEntity(age, name, npcId, int32(GParamCfg.GetEntityTypePerson()), 100),
	}

	WorldMap.GEntityList[uint32(npcId)] = e
	WorldMap.GEntityTypeList[uint32(npcId)] = GParamCfg.GetEntityTypePerson()

	e.AddActionLog(&ActionLog{
		ActionType: ActionParamCfg.ActionBorn,
		Action:     "你出生了,你的目标是尽可能在生命中做更多有意义的事情,以及活着((生命值,饱食度 > 0))",
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
	})

	e.RegisterAction()
	go e.LifeProcess()

	return e
}

func (u *EntityUser) RegisterAction() {
	u.Register(ActionParamCfg.ActionTypeMove, ActionMove)
	u.Register(ActionParamCfg.ActionTypeCuttingDownTrees, ActionCuttingDownTrees)

	u.Register(ActionParamCfg.ActionSleep, ActionSleep)

	u.Register(ActionParamCfg.ActionTypeBreakFirst, ActionEatDinner)
	u.Register(ActionParamCfg.ActionTypeLunch, ActionEatDinner)
	u.Register(ActionParamCfg.ActionDinner, ActionEatDinner)

}

func (u *EntityUser) LifeProcess() {
	ticker := time.NewTicker(1 * time.Second)
	tickerten := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			u.SatietyDegree -= 1
			if u.HP <= 0 {
				WorldMap.Gmap.DeadChan <- u
				return
			}
		case <-tickerten.C:
			u.Age++
			u.HP++
		}

	}
}
