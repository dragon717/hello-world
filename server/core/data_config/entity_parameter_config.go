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

type EntityParameterList struct {
	Key     string `xml:"key,attr"`
	Content string `xml:"content,attr"`
	Type    string
}

type EntityParameterConfig struct {
	PList []*EntityParameterList `xml:"data"`
}

type XmlEntityParameterConfig struct {
	EntityPerson uint32
	EntityTree   uint32
	EntityStone  uint32
}

func (this *XmlEntityParameterConfig) LoadConfig() bool {
	if err := this.LoadXmlEntityParameterConfig(); err != nil {
		return false
	}
	return true
}

func (data_sheet *XmlEntityParameterConfig) LoadXmlEntityParameterConfig(path ...string) error {
	var filename string
	if len(path) == 0 {
		filename = DataConfigDir + "entity_parameter.xml"
	} else {
		filename = path[0]
	}

	var pListData = &EntityParameterConfig{
		PList: make([]*EntityParameterList, 0),
	}
	err := common.LoadConfig(filename, &pListData)
	if err != nil {
		return err
	}

	var dataMap = make(map[string]*EntityParameterList)
	for _, val := range pListData.PList {
		dataMap[val.Key] = val
	}

	{
		var dataValStruct, exist = dataMap["entity_person"]
		if !exist {
			return fmt.Errorf("EntityPerson Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.EntityPerson = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["entity_tree"]
		if !exist {
			return fmt.Errorf("EntityTree Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.EntityTree = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["entity_stone"]
		if !exist {
			return fmt.Errorf("EntityStone Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.EntityStone = uint32(intVal)

	}

	return nil
}

func GetXmlEntityParameterConfig() *XmlEntityParameterConfig {
	return &XmlEntityParameterConfig{}
}

// 人类
func (this *XmlEntityParameterConfig) GetEntityPerson() uint32 {
	return this.EntityPerson
}

// 树
func (this *XmlEntityParameterConfig) GetEntityTree() uint32 {
	return this.EntityTree
}

// 巨石
func (this *XmlEntityParameterConfig) GetEntityStone() uint32 {
	return this.EntityStone
}
