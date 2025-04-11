package main

import (
	"fmt"
	"math/rand"
)

func ActionMove(op *ActionMsg, u EntityInterface) {
	WorldMap.Gmap.MoveLocation(int(u.GetX()), int(u.GetY()), op.Target.Position.X, op.Target.Position.Y, u)
	u.AddActionLog(&ActionLog{
		ActionType: op.Action,
		Action:     fmt.Sprintf("移动到[%d,%d]", op.Target.Position.X, op.Target.Position.Y),
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
	})
}

func ActionCuttingDownTrees(op *ActionMsg, u EntityInterface) {
	target := FindTarget(uint32(op.Target.Entity.ID))
	if target == nil || op.Target.Entity.ID == op.SelfID {
		fmt.Println("非法目标")
	}

	if u.GetBag()[ItemParamCfg.GetItemFellingaxe()] == 0 {
		u.AddActionLog(&ActionLog{
			ActionType: op.Action,
			Action:     op.Reason,
			Time:       WorldMap.Gmap.GlobalTime.GetTime(),
			Result:     "你没有斧头,无法砍树",
		})
		return
	}

	u.AddActionLog(&ActionLog{
		ActionType: op.Action,
		Action:     op.Reason,
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
	})
	target.SendActionChan(op)
}
func BeActionCuttingDownTrees(op *ActionMsg, u EntityInterface) {
	log := u.AddActionLog(&ActionLog{
		ActionType: op.Action,
		Action:     "被砍了",
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
	})
	res := &ResultMsg{
		Ret:      RETNO_OK,
		ActionID: log.ID,
		Result:   "砍伐成功!获得木材",
	}
	if u.GetHP() <= 0 {
		res.Result = "树已死亡,砍伐失败"
		res.Ret = RETNO_ERROR
	}
	u.SetHP(u.GetHP() - int32(rand.Intn(10)+10))
	res.AwardList = map[uint32]uint32{ItemParamCfg.GetItemWood(): uint32(rand.Intn(2) + 1)}

	fmt.Println("树", u.GetId(), "被砍了", "剩余HP:", u.GetHP())

	c := FindTarget(uint32(op.SelfID))
	c.SendCallBackChan(res)
}
func ActionEatDinner(op *ActionMsg, u EntityInterface) {
	log := &ActionLog{
		ActionType: op.Action,
		Action:     op.Reason,
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
	}
	if u.GetBag()[uint32(op.Target.Item.ItemID)] == 0 {
		log.Result = "你的背包里没有这个食物"
		u.AddActionLog(log)
		return
	}
	num := u.GetBag()[uint32(op.Target.Item.ItemID)]
	u.SetBagItem(uint32(op.Target.Item.ItemID), num-uint32(op.Target.Item.Count))

	r := rand.Intn(10) + 20*int(num)
	if u.GetSatietyDegree()+int32(r) > 100 {
		u.SetSatietyDegree(100)
	} else {
		u.SetSatietyDegree(u.GetSatietyDegree() + int32(r))
	}
	log.Result = fmt.Sprintf("食用%d个,饱食度+%d", op.Target.Item.Count, r)

	u.AddActionLog(log)
}
func ActionSleep(op *ActionMsg, u EntityInterface) {
	log := &ActionLog{
		ActionType: op.Action,
		Action:     op.Reason,
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
	}
	if u.GetSatietyDegree() < 10 {
		log.Result = "你太饥饿了,不能睡觉"
		u.AddActionLog(log)
		return
	}
	r := rand.Intn(5) + 5
	if int(u.GetSatietyDegree())-r <= 0 {
		u.SetSatietyDegree(1)
	} else {
		u.SetSatietyDegree(u.GetSatietyDegree() - int32(r))
	}

	if u.GetHP()+10 > 100 {
		u.SetHP(100)
	} else {
		u.SetHP(u.GetHP() + 10)
	}
	log.Result = fmt.Sprintf("饱食度-%d HP+10", r)
	u.AddActionLog(log)
}
func ActionPickUp(op *ActionMsg, u EntityInterface) {
	it := ItemCfg.GetById(int(op.Target.Item.ItemID))
	if it == nil {
		u.AddActionLog(&ActionLog{
			ActionType: op.Action,
			Action:     op.Reason,
			Time:       WorldMap.Gmap.GlobalTime.GetTime(),
			Result:     "周围没有这个物品或参数有误",
		})
		return
	}
	u.AddBagItem(uint32(op.Target.Item.ItemID), uint32(op.Target.Item.Count))
	log := &ActionLog{
		ActionType: op.Action,
		Action:     op.Reason,
		Time:       WorldMap.Gmap.GlobalTime.GetTime(),
		Result:     fmt.Sprintf("获得%d个%s", op.Target.Item.Count, it.Name),
	}
	u.AddActionLog(log)
}
