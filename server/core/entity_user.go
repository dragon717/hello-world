package main

import (
	"math/rand"
	"time"
)

type EntityUser struct {
	Tool  []string
	Money uint64
	*Entity
}

func NewUser(name string, npcId int32, age uint32, x, y int) *EntityUser {
	e := &EntityUser{
		Tool:   make([]string, 0),
		Money:  0,
		Entity: NewEntity(age, name, npcId, int32(EntityParamCfg.GetEntityPerson()), 100, x, y),
	}

	e.AddActionLog(&ActionLog{
		ActionType: ActionParamCfg.ActionBorn,
		Action:     "尽可能在生命中做更多有意义的事情,以及活着((生命值,饱食度 > 0))",
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
	u.Register(ActionParamCfg.ActionTypePickUp, ActionPickUp)

}

func (u *EntityUser) LifeProcess() {
	ticker := time.NewTicker(1 * time.Second)
	tickerTen := time.NewTicker(10 * time.Second)
	tickerThree := time.NewTicker(3 * ACTION_RATE * time.Second)
	for {
		select {
		case <-ticker.C:
			u.SatietyDegree -= 1
			if u.SatietyDegree <= 1 {
				u.HP -= 1
			}
			if u.HP <= 0 {
				WorldMap.Gmap.DeadChan <- u
				return
			}
		case <-tickerThree.C:
			_ = sendmsg(u)
			if len(u.GetActionLog())%10 == 0 {
				SendTargetTaskMsg(u)
			}
		case <-tickerTen.C:
			if u.Age > uint32(80+rand.Intn(10)) {
				WorldMap.Gmap.DeadChan <- u
				return
			}
			u.Age++
			u.HP++
		}

	}
}
