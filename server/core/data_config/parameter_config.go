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

type ParameterList struct {
	Key     string `xml:"key,attr"`
	Content string `xml:"content,attr"`
	Type    string
}

type ParameterConfig struct {
	PList []*ParameterList `xml:"data"`
}

type XmlParameterConfig struct {
	EntityTypePerson uint32
	EntityTypeTree   uint32
	ToolTypeAxe      uint32
	ItemTypeWood     uint32
}

func (this *XmlParameterConfig) LoadConfig() bool {
	if err := this.LoadXmlParameterConfig(); err != nil {
		return false
	}
	return true
}

func (data_sheet *XmlParameterConfig) LoadXmlParameterConfig(path ...string) error {
	var filename string
	if len(path) == 0 {
		filename = DataConfigDir + "parameter.xml"
	} else {
		filename = path[0]
	}

	var pListData = &ParameterConfig{
		PList: make([]*ParameterList, 0),
	}
	err := common.LoadConfig(filename, &pListData)
	if err != nil {
		return err
	}

	var dataMap = make(map[string]*ParameterList)
	for _, val := range pListData.PList {
		dataMap[val.Key] = val
	}

	{
		var dataValStruct, exist = dataMap["entity_type_person"]
		if !exist {
			return fmt.Errorf("EntityTypePerson Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.EntityTypePerson = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["entity_type_tree"]
		if !exist {
			return fmt.Errorf("EntityTypeTree Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.EntityTypeTree = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["tool_type_axe"]
		if !exist {
			return fmt.Errorf("ToolTypeAxe Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ToolTypeAxe = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["item_type_wood"]
		if !exist {
			return fmt.Errorf("ItemTypeWood Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ItemTypeWood = uint32(intVal)

	}

	return nil
}

func GetXmlParameterConfig() *XmlParameterConfig {
	return &XmlParameterConfig{}
}

// 人类
func (this *XmlParameterConfig) GetEntityTypePerson() uint32 {
	return this.EntityTypePerson
}

// 树
func (this *XmlParameterConfig) GetEntityTypeTree() uint32 {
	return this.EntityTypeTree
}

// 斧头
func (this *XmlParameterConfig) GetToolTypeAxe() uint32 {
	return this.ToolTypeAxe
}

// 木头
func (this *XmlParameterConfig) GetItemTypeWood() uint32 {
	return this.ItemTypeWood
}
