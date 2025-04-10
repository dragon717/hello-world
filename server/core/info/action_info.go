package info

import "Test/data_config"

type ActionInfo struct {
	List []*data_config.ActionInfo
	Map  map[int]*data_config.ActionInfo
}

func LoadActionCfg() *ActionInfo {
	e := &ActionInfo{
		List: make([]*data_config.ActionInfo, 0),
		Map:  make(map[int]*data_config.ActionInfo),
	}
	e.List = data_config.GetXmlActionInfo().Datas
	for _, info := range e.List {
		e.Map[info.Id] = info
	}
	return e
}
func (e *ActionInfo) GetById(id int) *data_config.ActionInfo {
	return e.Map[id]
}
func (e *ActionInfo) GetByName(name string) *data_config.ActionInfo {
	for _, info := range e.List {
		if info.Name == name {
			return info
		}
	}
	return nil
}
