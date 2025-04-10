package info

import "Test/data_config"

type EntityInfo struct {
	List []*data_config.EntityInfo
	Map  map[int]*data_config.EntityInfo
}

func LoadEntityCfg() *EntityInfo {
	e := &EntityInfo{
		List: make([]*data_config.EntityInfo, 0),
		Map:  make(map[int]*data_config.EntityInfo),
	}
	e.List = data_config.GetXmlEntityInfo().Datas
	for _, info := range e.List {
		e.Map[info.Id] = info
	}
	return e
}
func (e *EntityInfo) GetById(id int) *data_config.EntityInfo {
	return e.Map[id]
}
func (e *EntityInfo) GetByName(name string) *data_config.EntityInfo {
	for _, info := range e.List {
		if info.Name == name {
			return info
		}
	}
	return nil
}
