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

type ItemTypeParameterList struct {
	Key     string `xml:"key,attr"`
	Content string `xml:"content,attr"`
	Type    string
}

type ItemTypeParameterConfig struct {
	PList []*ItemTypeParameterList `xml:"data"`
}

type XmlItemTypeParameterConfig struct {
	ItemTypeTool     uint32
	ItemTypeFood     uint32
	ItemTypeToy      uint32
	ItemTypeMaterial uint32
}

func (this *XmlItemTypeParameterConfig) LoadConfig() bool {
	if err := this.LoadXmlItemTypeParameterConfig(); err != nil {
		return false
	}
	return true
}

func (data_sheet *XmlItemTypeParameterConfig) LoadXmlItemTypeParameterConfig(path ...string) error {
	var filename string
	if len(path) == 0 {
		filename = DataConfigDir + "item_type_parameter.xml"
	} else {
		filename = path[0]
	}

	var pListData = &ItemTypeParameterConfig{
		PList: make([]*ItemTypeParameterList, 0),
	}
	err := common.LoadConfig(filename, &pListData)
	if err != nil {
		return err
	}

	var dataMap = make(map[string]*ItemTypeParameterList)
	for _, val := range pListData.PList {
		dataMap[val.Key] = val
	}

	{
		var dataValStruct, exist = dataMap["item_type_tool"]
		if !exist {
			return fmt.Errorf("ItemTypeTool Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ItemTypeTool = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["item_type_food"]
		if !exist {
			return fmt.Errorf("ItemTypeFood Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ItemTypeFood = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["item_type_toy"]
		if !exist {
			return fmt.Errorf("ItemTypeToy Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ItemTypeToy = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["item_type_material"]
		if !exist {
			return fmt.Errorf("ItemTypeMaterial Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ItemTypeMaterial = uint32(intVal)

	}

	return nil
}

func GetXmlItemTypeParameterConfig() *XmlItemTypeParameterConfig {
	return &XmlItemTypeParameterConfig{}
}

// 工具
func (this *XmlItemTypeParameterConfig) GetItemTypeTool() uint32 {
	return this.ItemTypeTool
}

// 食物
func (this *XmlItemTypeParameterConfig) GetItemTypeFood() uint32 {
	return this.ItemTypeFood
}

// 玩具
func (this *XmlItemTypeParameterConfig) GetItemTypeToy() uint32 {
	return this.ItemTypeToy
}

// 材料
func (this *XmlItemTypeParameterConfig) GetItemTypeMaterial() uint32 {
	return this.ItemTypeMaterial
}
