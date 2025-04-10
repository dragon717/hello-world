package data_config

// 该文件为工具生成 不要修改
import (
	"Test/common"
	"fmt"
	"strconv"
	"strings"
)

var _ = strings.Split("1|2|3", "|")
var _, _ = strconv.Atoi("1")

type ItemParameterList struct {
	Key     string `xml:"key,attr"`
	Content string `xml:"content,attr"`
	Type    string
}

type ItemParameterConfig struct {
	PList []*ItemParameterList `xml:"data"`
}

type XmlItemParameterConfig struct {
	ItemFellingaxe uint32
	ItemPickaxe    uint32
	ItemBerry      uint32
	ItemWood       uint32
	ItemBranch     uint32
	ItemPebble     uint32
}

func (this *XmlItemParameterConfig) LoadConfig() bool {
	if err := this.LoadXmlItemParameterConfig(); err != nil {
		return false
	}
	return true
}

func (data_sheet *XmlItemParameterConfig) LoadXmlItemParameterConfig(path ...string) error {
	var filename string
	if len(path) == 0 {
		filename = DataConfigDir + "item_parameter.xml"
	} else {
		filename = path[0]
	}

	var pListData = &ItemParameterConfig{
		PList: make([]*ItemParameterList, 0),
	}
	err := common.LoadConfig(filename, &pListData)
	if err != nil {
		return err
	}

	var dataMap = make(map[string]*ItemParameterList)
	for _, val := range pListData.PList {
		dataMap[val.Key] = val
	}

	{
		var dataValStruct, exist = dataMap["item_FellingAxe"]
		if !exist {
			return fmt.Errorf("ItemFellingaxe Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ItemFellingaxe = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["item_Pickaxe"]
		if !exist {
			return fmt.Errorf("ItemPickaxe Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ItemPickaxe = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["item_berry"]
		if !exist {
			return fmt.Errorf("ItemBerry Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ItemBerry = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["item_wood"]
		if !exist {
			return fmt.Errorf("ItemWood Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ItemWood = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["item_branch"]
		if !exist {
			return fmt.Errorf("ItemBranch Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ItemBranch = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["item_pebble"]
		if !exist {
			return fmt.Errorf("ItemPebble Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ItemPebble = uint32(intVal)

	}

	return nil
}

func GetXmlItemParameterConfig() *XmlItemParameterConfig {
	return &XmlItemParameterConfig{}
}

// 伐木斧(采集工具)
func (this *XmlItemParameterConfig) GetItemFellingaxe() uint32 {
	return this.ItemFellingaxe
}

// 采石镐(采集工具)
func (this *XmlItemParameterConfig) GetItemPickaxe() uint32 {
	return this.ItemPickaxe
}

// 浆果(回复饱食度)
func (this *XmlItemParameterConfig) GetItemBerry() uint32 {
	return this.ItemBerry
}

// 木头(建筑材料)
func (this *XmlItemParameterConfig) GetItemWood() uint32 {
	return this.ItemWood
}

// 树枝(工具材料)
func (this *XmlItemParameterConfig) GetItemBranch() uint32 {
	return this.ItemBranch
}

// 鹅卵石(工具材料)
func (this *XmlItemParameterConfig) GetItemPebble() uint32 {
	return this.ItemPebble
}
