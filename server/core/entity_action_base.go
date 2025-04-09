package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

func ActionMove(op *ActionMsg, u EntityInterface) {
	WorldMap.Gmap.MoveLocation(int(u.GetX()), int(u.GetY()), op.Target.Y, op.Target.X, u)
	u.AddActionLog(&ActionLog{
		ActionType: op.Action,
		Action:     fmt.Sprintf("移动到[%d,%d]", op.Target.X, op.Target.Y),
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
	})
}

func ActionCuttingDownTrees(op *ActionMsg, u EntityInterface) {
	u.AddActionLog(&ActionLog{
		ActionType: op.Action,
		Action:     "砍树",
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
	})
}
func BeActionCuttingDownTrees(op *ActionMsg, u EntityInterface) {
	if u.GetHP() <= 0 {
	}
	u.SetHP(u.GetHP() - int32(rand.Intn(10)+10))
	log := u.AddActionLog(&ActionLog{
		ActionType: op.Action,
		Action:     "被砍了",
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
	})
	fmt.Println("树", u.GetId(), "被砍了", "剩余HP:", u.GetHP())
	c := FindTarget(uint32(op.SelfId))
	c.SendCallBackChan(&ResultMsg{
		ActionID:  log.ID,
		Result:    "砍伐成功!获得木材",
		AwardList: map[uint32]uint32{GParamCfg.GetItemTypeWood(): uint32(rand.Intn(2) + 1)},
	})
}
func ActionEatDinner(op *ActionMsg, u EntityInterface) {
	r := rand.Intn(10) + 20
	if u.GetSatietyDegree()+uint32(r) > 100 {
		u.SetSatietyDegree(100)
	} else {
		u.SetSatietyDegree(u.GetSatietyDegree() + uint32(r))

	}

	u.AddActionLog(&ActionLog{
		ActionType: op.Action,
		Action:     op.Reason,
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
		Result:     "饱食度+" + strconv.Itoa(r),
	})
}
func ActionSleep(op *ActionMsg, u EntityInterface) {
	var Result string
	r := rand.Intn(10) + 10
	if int(u.GetSatietyDegree())-r <= 30 {
		u.SetSatietyDegree(1)
	} else {
		u.SetSatietyDegree(u.GetSatietyDegree() - uint32(r))
	}

	if u.GetHP()+10 > 100 {
		u.SetHP(100)
	} else {
		u.SetHP(u.GetHP() + 10)
	}
	Result += fmt.Sprintf("饱食度-%d HP+10", r)
	u.AddActionLog(&ActionLog{
		ActionType: op.Action,
		Action:     op.Reason,
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
		Result:     Result,
	})
}
