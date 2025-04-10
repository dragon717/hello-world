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

type ActionTypeParameterList struct {
	Key     string `xml:"key,attr"`
	Content string `xml:"content,attr"`
	Type    string
}

type ActionTypeParameterConfig struct {
	PList []*ActionTypeParameterList `xml:"data"`
}

type XmlActionTypeParameterConfig struct {
	ActionTypeInitiative uint32
	EntityTypePassive    uint32
	EntityTypeMutual     uint32
}

func (this *XmlActionTypeParameterConfig) LoadConfig() bool {
	if err := this.LoadXmlActionTypeParameterConfig(); err != nil {
		return false
	}
	return true
}

func (data_sheet *XmlActionTypeParameterConfig) LoadXmlActionTypeParameterConfig(path ...string) error {
	var filename string
	if len(path) == 0 {
		filename = DataConfigDir + "action_type_parameter.xml"
	} else {
		filename = path[0]
	}

	var pListData = &ActionTypeParameterConfig{
		PList: make([]*ActionTypeParameterList, 0),
	}
	err := common.LoadConfig(filename, &pListData)
	if err != nil {
		return err
	}

	var dataMap = make(map[string]*ActionTypeParameterList)
	for _, val := range pListData.PList {
		dataMap[val.Key] = val
	}

	{
		var dataValStruct, exist = dataMap["action_type_initiative"]
		if !exist {
			return fmt.Errorf("ActionTypeInitiative Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ActionTypeInitiative = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["entity_type_passive"]
		if !exist {
			return fmt.Errorf("EntityTypePassive Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.EntityTypePassive = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["entity_type_mutual"]
		if !exist {
			return fmt.Errorf("EntityTypeMutual Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.EntityTypeMutual = uint32(intVal)

	}

	return nil
}

func GetXmlActionTypeParameterConfig() *XmlActionTypeParameterConfig {
	return &XmlActionTypeParameterConfig{}
}

// 主动
func (this *XmlActionTypeParameterConfig) GetActionTypeInitiative() uint32 {
	return this.ActionTypeInitiative
}

// 被动
func (this *XmlActionTypeParameterConfig) GetEntityTypePassive() uint32 {
	return this.EntityTypePassive
}

// 共有
func (this *XmlActionTypeParameterConfig) GetEntityTypeMutual() uint32 {
	return this.EntityTypeMutual
}
