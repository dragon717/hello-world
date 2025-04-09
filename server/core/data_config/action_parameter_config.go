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

type ActionParameterList struct {
	Key     string `xml:"key,attr"`
	Content string `xml:"content,attr"`
	Type    string
}

type ActionParameterConfig struct {
	PList []*ActionParameterList `xml:"data"`
}

type XmlActionParameterConfig struct {
	ActionTypeMove             uint32
	ActionTypeCuttingDownTrees uint32
	ActionTypeBreakFirst       uint32
	ActionTypeLunch            uint32
	ActionDinner               uint32
	ActionSleep                uint32
	ActionGrow                 uint32
	ActionBorn                 uint32
}

func (this *XmlActionParameterConfig) LoadConfig() bool {
	if err := this.LoadXmlActionParameterConfig(); err != nil {
		return false
	}
	return true
}

func (data_sheet *XmlActionParameterConfig) LoadXmlActionParameterConfig(path ...string) error {
	var filename string
	if len(path) == 0 {
		filename = DataConfigDir + "action_parameter.xml"
	} else {
		filename = path[0]
	}

	var pListData = &ActionParameterConfig{
		PList: make([]*ActionParameterList, 0),
	}
	err := common.LoadConfig(filename, &pListData)
	if err != nil {
		return err
	}

	var dataMap = make(map[string]*ActionParameterList)
	for _, val := range pListData.PList {
		dataMap[val.Key] = val
	}

	{
		var dataValStruct, exist = dataMap["action_type_move"]
		if !exist {
			return fmt.Errorf("ActionTypeMove Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ActionTypeMove = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["action_type_Cutting_down_trees"]
		if !exist {
			return fmt.Errorf("ActionTypeCuttingDownTrees Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ActionTypeCuttingDownTrees = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["action_type_break_first"]
		if !exist {
			return fmt.Errorf("ActionTypeBreakFirst Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ActionTypeBreakFirst = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["action_type_lunch"]
		if !exist {
			return fmt.Errorf("ActionTypeLunch Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ActionTypeLunch = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["action_dinner"]
		if !exist {
			return fmt.Errorf("ActionDinner Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ActionDinner = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["action_sleep"]
		if !exist {
			return fmt.Errorf("ActionSleep Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ActionSleep = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["action_grow"]
		if !exist {
			return fmt.Errorf("ActionGrow Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ActionGrow = uint32(intVal)

	}

	{
		var dataValStruct, exist = dataMap["action_born"]
		if !exist {
			return fmt.Errorf("ActionBorn Not Exist")
		}
		var dataVal = dataValStruct.Content

		intVal, err := strconv.ParseInt(dataVal, 10, 64)
		if err != nil {
			return err
		}
		data_sheet.ActionBorn = uint32(intVal)

	}

	return nil
}

func GetXmlActionParameterConfig() *XmlActionParameterConfig {
	return &XmlActionParameterConfig{}
}

// 移动
func (this *XmlActionParameterConfig) GetActionTypeMove() uint32 {
	return this.ActionTypeMove
}

// 砍树
func (this *XmlActionParameterConfig) GetActionTypeCuttingDownTrees() uint32 {
	return this.ActionTypeCuttingDownTrees
}

// 吃早餐
func (this *XmlActionParameterConfig) GetActionTypeBreakFirst() uint32 {
	return this.ActionTypeBreakFirst
}

// 吃午饭
func (this *XmlActionParameterConfig) GetActionTypeLunch() uint32 {
	return this.ActionTypeLunch
}

// 吃晚饭
func (this *XmlActionParameterConfig) GetActionDinner() uint32 {
	return this.ActionDinner
}

// 睡觉
func (this *XmlActionParameterConfig) GetActionSleep() uint32 {
	return this.ActionSleep
}

// 生长
func (this *XmlActionParameterConfig) GetActionGrow() uint32 {
	return this.ActionGrow
}

// 出生
func (this *XmlActionParameterConfig) GetActionBorn() uint32 {
	return this.ActionBorn
}
