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

type EntityTypeParameterList struct {
	Key     string `xml:"key,attr"`
	Content string `xml:"content,attr"`
	Type    string
}

type EntityTypeParameterConfig struct {
	PList []*EntityTypeParameterList `xml:"data"`
}

type XmlEntityTypeParameterConfig struct {
	EntityTypeAnimal      uint32
	EntityTypePlant       uint32
	EntityTypeBuilding    uint32
	EntityTypeEnvironment uint32
}

func (this *XmlEntityTypeParameterConfig) LoadConfig() bool {
	if err := this.LoadXmlEntityTypeParameterConfig(); err != nil {
		return false
	}
	return true
}

func (data_sheet *XmlEntityTypeParameterConfig) LoadXmlEntityTypeParameterConfig(path ...string) error {
	var filename string
	if len(path) == 0 {
		filename = DataConfigDir + "entity_type_parameter.xml"
	} else {
		filename = path[0]
	}

	var pListData = &EntityTypeParameterConfig{
		PList: make([]*EntityTypeParameterList, 0),
	}
	err := common.LoadConfig(filename, &pListData)
	if err != nil {
		return err
	}

	var dataMap = make(map[string]*EntityTypeParameterList)
	for _, val := range pListData.PList {
		dataMap[val.Key] = val
	}

	{
		var dataValStruct, exist = dataMap["entity_type_animal"]
		if !exist {
			return fmt.Errorf("EntityTypeAnimal Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.EntityTypeAnimal = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["entity_type_plant"]
		if !exist {
			return fmt.Errorf("EntityTypePlant Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.EntityTypePlant = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["entity_type_building"]
		if !exist {
			return fmt.Errorf("EntityTypeBuilding Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.EntityTypeBuilding = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["entity_type_environment"]
		if !exist {
			return fmt.Errorf("EntityTypeEnvironment Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.EntityTypeEnvironment = uint32(intVal)

	}

	return nil
}

func GetXmlEntityTypeParameterConfig() *XmlEntityTypeParameterConfig {
	return &XmlEntityTypeParameterConfig{}
}

// 动物
func (this *XmlEntityTypeParameterConfig) GetEntityTypeAnimal() uint32 {
	return this.EntityTypeAnimal
}

// 植物
func (this *XmlEntityTypeParameterConfig) GetEntityTypePlant() uint32 {
	return this.EntityTypePlant
}

// 建筑
func (this *XmlEntityTypeParameterConfig) GetEntityTypeBuilding() uint32 {
	return this.EntityTypeBuilding
}

// 环境
func (this *XmlEntityTypeParameterConfig) GetEntityTypeEnvironment() uint32 {
	return this.EntityTypeEnvironment
}
