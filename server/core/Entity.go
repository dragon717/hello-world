package main

import (
	"fmt"
	"slices"
)

type EntityInterface interface {
	GetId() int32
	GetName() string
	GetType() int32
	GetHP() int32
	GetStatus() bool
	GetX() uint32
	GetY() uint32
	GetActionLog() []*ActionLog
	GetAge() uint32
	GetSatietyDegree() uint32
	GetBag() map[uint32]uint32

	SetX(x uint32)
	SetY(y uint32)
	SetStatus(status bool)
	SetHP(hp int32)
	SetId(id int32)
	SetName(name string)
	SetType(ty int32)
	SetAge(age uint32)
	SetSatietyDegree(satietyDegree uint32)

	SetBagItem(itemId uint32, num uint32)
	AddBagItem(itemId uint32, num uint32)

	GetInfo() map[string]string

	RegisterAction()

	AddActionLog(log *ActionLog) *ActionLog
	SendActionChan(op *ActionMsg)
	SendCallBackChan(op *ResultMsg)
	Register(uint32, func(*ActionMsg, EntityInterface))
}

type Entity struct {
	Id            int32
	Name          string
	Type          int32
	HP            int32
	Age           uint32
	X             uint32
	Y             uint32
	SatietyDegree uint32 //饱食度
	Bag           map[uint32]uint32
	Status        bool //是否可交互
	actionLog     []*ActionLog
	ResultChan    chan *ResultMsg //接收更新记忆的消息
	ActionChan    chan *ActionMsg //接收model消息
	ActionList    map[uint32]func(*ActionMsg, EntityInterface)
}

type ActionLog struct {
	ID         uint32
	ActionType uint32
	Action     string
	Time       string
	Result     string
}

func NewEntity(age uint32, name string, id, ty, hp int32, x, y int) *Entity {
	e := &Entity{
		Id:            id,
		Name:          name,
		Type:          ty,
		HP:            hp,
		X:             uint32(x),
		Y:             uint32(y),
		Age:           age,
		Bag:           make(map[uint32]uint32),
		SatietyDegree: 100,
		Status:        true,
		actionLog:     make([]*ActionLog, 0),
		ResultChan:    make(chan *ResultMsg, 1),
		ActionChan:    make(chan *ActionMsg, 1),
		ActionList:    make(map[uint32]func(*ActionMsg, EntityInterface)),
	}

	go e.ConsumerChan()
	return e
}

func FindTarget(ID uint32) EntityInterface {
	if target, ok := WorldMap.GEntityList[ID]; ok {
		return target
	}
	return nil
}

func (e *Entity) GetId() int32 {
	return e.Id
}

func (e *Entity) GetName() string {
	return e.Name
}

func (e *Entity) GetType() int32 {
	return e.Type
}
func (e *Entity) GetHP() int32 {
	return e.HP
}
func (e *Entity) GetX() uint32 {
	return e.X
}
func (e *Entity) GetY() uint32 {
	return e.Y
}
func (e *Entity) GetStatus() bool {
	return e.Status
}

func (e *Entity) GetAge() uint32 {
	return e.Age
}
func (e *Entity) GetActionLog() []*ActionLog {
	return e.actionLog
}
func (e *Entity) GetSatietyDegree() uint32 {
	return e.SatietyDegree
}
func (e *Entity) SetBagItem(itemId uint32, num uint32) {
	e.Bag[itemId] = num
}
func (e *Entity) AddBagItem(itemId uint32, num uint32) {
	e.Bag[itemId] += num
}
func (e *Entity) GetBag() map[uint32]uint32 {
	return e.Bag
}
func (e *Entity) SetStatus(status bool) {
	e.Status = status
}
func (e *Entity) SetHP(hp int32) {
	e.HP = hp
}
func (e *Entity) SetId(id int32) {
	e.Id = id
}
func (e *Entity) SetName(name string) {
	e.Name = name
}
func (e *Entity) SetType(ty int32) {
	e.Type = ty
}
func (e *Entity) SetX(x uint32) {
	e.X = x
}
func (e *Entity) SetY(y uint32) {
	e.Y = y
}
func (e *Entity) SetAge(age uint32) {
	e.Age = age
}
func (e *Entity) SetSatietyDegree(satietyDegree uint32) {
	e.SatietyDegree = satietyDegree
}
func (e *Entity) AddActionLog(log *ActionLog) *ActionLog {
	log.ID = uint32(len(e.actionLog) + 1)
	e.actionLog = append(e.actionLog, log)
	return log
}
func (e *Entity) Register(id uint32, f func(*ActionMsg, EntityInterface)) {
	e.ActionList[id] = f
}
func (e *Entity) SendActionChan(op *ActionMsg) {
	e.ActionChan <- op
}
func (e *Entity) SendCallBackChan(op *ResultMsg) {
	e.ResultChan <- op
}
func (e *Entity) ConsumerChan() {
	for {
		select {
		case op := <-e.ResultChan:
			index := slices.IndexFunc(e.actionLog, func(log *ActionLog) bool {
				return log.ID == op.ActionID
			})
			if index == -1 {
				continue
			}
			e.actionLog[index].Result = op.Result
			for ty, num := range op.AwardList {
				e.Bag[ty] += num
			}
			fmt.Println(op.ActionID)

		case op := <-e.ActionChan:
			if f, ok := e.ActionList[op.Action]; ok {
				f(op, e)
			} else {
				fmt.Println(e.Name, "非法操作", op.Action)
			}
		}
	}
}

func (e *Entity) RegisterAction() {
}

func (e *Entity) GetInfo() map[string]string {
	return map[string]string{
		"id":            fmt.Sprintf("%d", e.Id),
		"name":          e.Name,
		"type":          fmt.Sprintf("%d", e.Type),
		"hp":            fmt.Sprintf("%d", e.HP),
		"age":           fmt.Sprintf("%d", e.Age),
		"x":             fmt.Sprintf("%d", e.X),
		"y":             fmt.Sprintf("%d", e.Y),
		"satietydegree": fmt.Sprintf("%d", e.SatietyDegree),
		"status":        fmt.Sprintf("%t", e.Status),
		"bag":           fmt.Sprintf("%v", e.Bag),
	}
}
