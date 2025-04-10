package main

import (
	"Test/data_config"
	"Test/info"
	"time"
)

// cfg:具体数据
// ParamCfg :函数式获取ID(避免静态变量)
// TypeParamCfg:数据类型

var ActionCfg = info.LoadActionCfg()
var ActionParamCfg *data_config.XmlActionParameterConfig
var ActionTypeParamCfg *data_config.XmlActionTypeParameterConfig

var EntityCfg = info.LoadEntityCfg()
var EntityParamCfg *data_config.XmlEntityParameterConfig
var EntityTypeParamCfg *data_config.XmlEntityTypeParameterConfig

var ItemCfg = info.LoadItemCfg()
var ItemParamCfg *data_config.XmlItemParameterConfig
var ItemtypeParamCfg *data_config.XmlItemTypeParameterConfig

var ACTION_RATE time.Duration = 2 //api请求速率

func InitCfg() {

	ActionParamCfg = data_config.GetXmlActionParameterConfig()
	ActionTypeParamCfg = data_config.GetXmlActionTypeParameterConfig()

	EntityParamCfg = data_config.GetXmlEntityParameterConfig()
	EntityTypeParamCfg = data_config.GetXmlEntityTypeParameterConfig()
	ItemParamCfg = data_config.GetXmlItemParameterConfig()
	ItemtypeParamCfg = data_config.GetXmlItemTypeParameterConfig()

	ActionParamCfg.LoadConfig()
	ActionTypeParamCfg.LoadConfig()

	EntityParamCfg.LoadConfig()
	EntityTypeParamCfg.LoadConfig()
	ItemParamCfg.LoadConfig()
	ItemtypeParamCfg.LoadConfig()
}
