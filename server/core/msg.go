package main

import (
	"Test/common"
	"fmt"
)

type ActionMsg struct {
	SelfId int32   `json:"selfid"`
	Action uint32  `json:"action"`
	Target *Target `json:"target"`
	Reason string  `json:"reason"`
}
type Target struct {
	X  int   `json:"x"`
	Y  int   `json:"y"`
	Id int32 `json:"id"`
}

type ResultMsg struct {
	ActionID  uint32 `json:"action"`
	Result    string `json:"result"`
	AwardList map[uint32]uint32
}

func formatActionMsg(you EntityInterface) string {
	var msg string

	var record string
	for _, actionLog := range you.GetActionLog() {
		record += fmt.Sprintf("%s %s %s,", actionLog.Time, actionLog.Action, actionLog.Result)
	}
	//record, _ := json.Marshal(you.GetActionLog())
	baseinfo := fmt.Sprintf("你的ID:%d;名字:%s;类型ID:%d;背包:%+v;坐标:[%d,%d],生命值:%d;饱食度:%d;记忆:(%s) ", you.GetId(), you.GetName(), you.GetType(), you.GetBag(), you.GetX(), you.GetY(), you.GetHP(), you.GetSatietyDegree(), string(record))
	mapInfo := fmt.Sprintf("当前时间:%s,地图周围信息:%s", WorldMap.Gmap.GlobalTime.GetTime(), WorldMap.Gmap.GetAroundInfo(int(you.GetX()), int(you.GetY()), 1))
	actionInfo := fmt.Sprintf("行为类型:%s", common.JsonMarshal(ActionParamCfg))
	entityInfo := fmt.Sprintf("实体类型:%s", common.JsonMarshal(EntityCfg.List))
	msg = fmt.Sprintf("%s,%s,%s,%s,你可以使用'行为类型'进行交互或者活动,你给出一个你想进行的行为,尽量不要重复,请用JSON格式响应以下请求，不要包含任何解释或格式化字符。JSON格式:{selfid:number,action:ACTION_TYPE,target:{x:number,y:number,id:number},reason:string},例子:{selfid:1001,action:100001,target:{x:0,y:2,id:8888},reason:''} 解释:你使用action:100001(ActionTypeCuttingDownTrees)对坐标[0,2]的id为8888的树进行了砍伐", baseinfo, mapInfo, actionInfo, entityInfo)
	return msg
}
