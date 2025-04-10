package main

type World struct {
	GcallBack       chan *ActionMsg //更新记忆
	Gmap            *Wmap
	GEntityTypeList map[uint32]uint32          //key:实体id value:实体类型
	GEntityList     map[uint32]EntityInterface //key:实体id value:实体
	Gpool           *Pool
}

func NewWorld() *World {
	return &World{
		GcallBack:       make(chan *ActionMsg),
		Gmap:            NewMap(8),
		GEntityTypeList: make(map[uint32]uint32),
		GEntityList:     make(map[uint32]EntityInterface),
		Gpool:           NewPool(100),
	}
}
func (w *World) AddEntity(entity EntityInterface) {
	w.GEntityList[uint32(entity.GetId())] = entity
	w.GEntityTypeList[uint32(entity.GetId())] = uint32(entity.GetType())
	w.Gmap.SetLocation(int(entity.GetX()), int(entity.GetY()), entity)
}
