package main

import (
	"Test/common"
	"Test/protocol/retno"
	"fmt"
)

// 统一行为结构体
type ActionMsg struct {
	SelfID int32        `json:"selfid"`
	Action uint32       `json:"action"` //
	Target ActionTarget `json:"target"` // 通用目标容器
	Reason string       `json:"reason"`
}

// 使用匿名结构体实现多态目标
type ActionTarget struct {
	Entity   *EntityTarget   `json:"entity,omitempty"`    // 实体目标
	Item     *ItemTarget     `json:"item,omitempty"`      // 物品目标
	Craft    *CraftTarget    `json:"craft,omitempty" `    // 合成操作
	Duration *DurationTarget `json:"duration,omitempty" ` // 时长相关
	Position *Position       `json:"position,omitempty"`  // 纯坐标目标
}

// 子目标类型定义
type EntityTarget struct {
	ID int32 `json:"id"` // 实体ID
	Position
}

type ItemTarget struct {
	ItemType int32 `json:"Item_type"` // 物品类型
	ItemID   int32 `json:"item_id"`   // 物品ID
	Count    int   `json:"count"`     // 操作数量
}

type CraftTarget struct {
	ItemType int32 `json:"item_type"` // 被合成的(物品/实体)类型
	ItemID   int32 `json:"item_id"`   // 被合成的(物品/实体)ID
}

type DurationTarget struct {
	time uint32 `json:"time"` // 持续时间（小时）
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ResultMsg struct {
	ActionID  uint32    `json:"action"`
	Ret       retno.RET `json:"ret"`
	Result    string    `json:"result"`
	AwardList map[uint32]uint32
}

func formatActionMsg(you EntityInterface) string {
	var msg string

	var record string
	for _, actionLog := range you.GetActionLog() {
		record += fmt.Sprintf("%s你%s", actionLog.Time, actionLog.Action)
		if len(actionLog.Result) != 0 {
			record += fmt.Sprintf("结果：%s,", actionLog.Result)
		}
	}
	//record, _ := json.Marshal(you.GetActionLog())
	baseinfo := fmt.Sprintf("你的名字是:%s(ID:%d;);类型ID:%d;背包:%+v;坐标:[%d,%d],当前任务目标:%s,生命值:%d;饱食度:%d;记忆:(%s) \n", you.GetName(), you.GetId(), you.GetType(), you.GetBag(), you.GetX(), you.GetY(), you.GetTargetTask(), you.GetHP(), you.GetSatietyDegree(), string(record))
	mapInfo := fmt.Sprintf("当前时间:%s,地图周围信息:%s\n", WorldMap.Gmap.GlobalTime.GetTime(), WorldMap.Gmap.GetAroundInfo(int(you.GetX()), int(you.GetY()), 1, you.GetId()))

	actionInfo := fmt.Sprintf("可用行为类型:%s\n", you.GetActionListInfo())
	entityInfo := fmt.Sprintf("实体大全:%s\n", common.JsonMarshal(EntityCfg.List))
	itemInfo := fmt.Sprintf("物品大全:%s\n", common.JsonMarshal(ItemCfg.List))

	if !(*devMode) {
		fmt.Print(baseinfo)
		fmt.Print(mapInfo)
		fmt.Print(actionInfo)
		fmt.Print(entityInfo)
		fmt.Print(itemInfo)
	}
	//entityInfo := fmt.Sprintf("实体类型:%s\n", common.JsonMarshal(EntityCfg.List))
	msg = fmt.Sprintf("%s%s%s%s%s你可以使用'行为类型'进行交互或者活动,请结合记忆(上下文)给出一个行为,行为的结果可能失败或者有前提条件,请用JSON格式响应以下请求,不要包含任何解释或格式化字符.JSON格式对应结构体:%s.用例:%s,%s,%s,%s,%s,%s", baseinfo, mapInfo, actionInfo, entityInfo, itemInfo, stu, example1, example2, example3, example4, example5, example6)
	return msg
}

type TaskMsg struct {
	Task     string `json:"task"`
	IsChange bool   `json:"change"`
}

func formatTargetTaskMsg(you EntityInterface) string {

	var msgstr string
	var record string
	for _, actionLog := range you.GetActionLog() {
		record += fmt.Sprintf("%s你%s", actionLog.Time, actionLog.Action)
		if len(actionLog.Result) != 0 {
			record += fmt.Sprintf("结果：%s,", actionLog.Result)
		}
	}
	baseinfo := fmt.Sprintf("你的名字是:%s(ID:%d;);类型ID:%d;背包:%+v;坐标:[%d,%d],当前任务目标:%s,生命值:%d;饱食度:%d;记忆:(%s) \n", you.GetName(), you.GetId(), you.GetType(), you.GetBag(), you.GetX(), you.GetY(), you.GetTargetTask(), you.GetHP(), you.GetSatietyDegree(), string(record))
	mapInfo := fmt.Sprintf("当前时间:%s,地图周围信息:%s\n", WorldMap.Gmap.GlobalTime.GetTime(), WorldMap.Gmap.GetAroundInfo(int(you.GetX()), int(you.GetY()), 1, you.GetId()))
	actionInfo := fmt.Sprintf("可用行为类型:%s\n", you.GetActionListInfo())
	entityInfo := fmt.Sprintf("实体大全:%s\n", common.JsonMarshal(EntityCfg.List))
	itemInfo := fmt.Sprintf("物品大全:%s\n", common.JsonMarshal(ItemCfg.List))
	msgstr = fmt.Sprintf("%s%s%s%s%s根据以上内容判断我是否完成了当前任务目标,如果完成了就给出一个新的任务目标,请用{task:STRING,change:BOOL}格式的json响应请求,不要包含任何解释或格式化字符.", baseinfo, mapInfo, actionInfo, entityInfo, itemInfo)

	return msgstr
}

var stu = `
// 统一行为结构体
type ActionMsg struct {
	SelfID int32        json:"selfid"
	Type   uint32       json:"action" // 使用新枚举类型
	Target ActionTarget json:"target" // 通用目标容器
	Reason string       json:"reason"
}
// 使用匿名结构体实现多态目标
type ActionTarget struct {
	Entity   *EntityTarget   json:"entity,omitempty"   // 实体目标
	Item     *ItemTarget     json:"item,omitempty"     // 物品目标
	Craft    *CraftTarget    json:"craft,omitempty"    // 合成操作
	Duration *DurationTarget json:"duration,omitempty" // 时长相关
	Position *Position       json:"position,omitempty" // 纯坐标目标
}
// 子目标类型定义
type EntityTarget struct {
	ID int32 json:"id" // 实体ID
	Position
}
type ItemTarget struct {
	ItemType int32 json:"slot"  // 物品类型
	ItemID   int32 json:"item"  // 物品ID
	Quantity int   json:"count" // 操作数量
}
type CraftTarget struct {
	ItemID int32 json:"recipe" // 被合成的(物品/实体)ID
}
type DurationTarget struct {
	Seconds uint32 json:"sec" // 持续时间（小时）
}
type Position struct {
	X int json:"x"
	Y int json:"y"
}`
var example1 = `
//移动操作（目标实体）
{
  "selfid": 1000,
  "action": 100060,
  "reason": "砍树",
  "target": {
    "entity": {
      "id": 2000,
      "x": 1,
      "y": 2
    }
  }
}`

var example2 = `
//物品操作（背包物品）
{
  "selfid": 782,
  "action": 100020,
  "reason": "食用浆果",
  "target": {
    "item": {
      "item_type": 2000,
      "item_id": 20001,
      "count": 1
    }
  }
}`
var example3 = `
// 合成操作
{
  "selfid": 223,
  "action": 100020,
  "reason": "合成伐木斧",
  "target": {
    "craft": {
      "Item_type": 1000,
      "Item_id": 10001,
    }
  }
}`
var example4 = `
 //持续操作（休息,睡觉）
{
  "selfid": 1000,
  "action": 100050,
  "reason": "休息",
  "target": {
    "duration": {
      "time": 2,
    }
  }
}`
var example5 = `
纯坐标操作（移动到位置）
json
{
  "selfid": 112,
  "action": 100010,
  "reason": "移动",
  "target": {
    "position": {
      "x": 45,
      "y": 23
    }
  }
}`
var example6 = `
复合操作（同时包含实体和物品）
json
{
  "selfid": 1000,
  "action": 100060,
  "reason": "用钥匙开门",
  "target": {
    "entity": {
      "id": 2000,
      "x": 10,
      "y": 20
    },
    "item": {
      "Item_type": 1000,
      "Item_id": 10001,
    }
  }
}`
