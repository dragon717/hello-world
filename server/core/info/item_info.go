package info

import "Test/data_config"

type ItemInfo struct {
	List []*data_config.ItemInfo
	Map  map[int]*data_config.ItemInfo
}

func LoadItemCfg() *ItemInfo {
	e := &ItemInfo{
		List: make([]*data_config.ItemInfo, 0),
		Map:  make(map[int]*data_config.ItemInfo),
	}
	e.List = data_config.GetXmlItemInfo().Datas
	for _, info := range e.List {
		e.Map[info.Id] = info
	}
	return e
}
func (e *ItemInfo) GetById(id int) *data_config.ItemInfo {
	return e.Map[id]
}
func (e *ItemInfo) GetByName(name string) *data_config.ItemInfo {
	for _, info := range e.List {
		if info.Name == name {
			return info
		}
	}
	return nil
}
