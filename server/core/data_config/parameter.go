// 该文件为工具生成 不要修改
package data_config

import (
	"Test/common"
	"unsafe"
)

func GetXmlParameterInfo() *XmlParameterInfo {
	data_sheet, err := LoadXmlParameterInfo()
	if err != nil {
		return nil
	}
	return data_sheet
}

type XmlParameterInfo struct {
	Datas []*ParameterInfo
}

func GetXmlParameterInfoName() string {
	return DataConfigDir + "parameter.xml"
}
func LoadXmlParameterInfo() (*XmlParameterInfo, error) {
	return LoadXmlParameterInfoByFileName(DataConfigDir + "parameter.xml")
}

func LoadXmlParameterInfoByFileName(filename string) (*XmlParameterInfo, error) {
	data_sheet := &struct {
		Datas []*struct {
			Key         string `xml:"key,attr"`
			Content     string `xml:"content,attr"`
			Ty          string `xml:"ty,attr"`
			Description string `xml:"description,attr"`
		} `xml:"data"`
	}{}

	if err := common.LoadConfig(filename, data_sheet); err != nil {
		return nil, err
	}
	return (*XmlParameterInfo)(unsafe.Pointer(data_sheet)), nil
}

type ParameterInfo struct {
	Key         string // 序号
	Content     string // 参数内容
	Ty          string // 参数类型
	Description string // 细节描述
}

func (this *ParameterInfo) GetKey() string {
	return this.Key
}

func (this *ParameterInfo) GetContent() string {
	return this.Content
}

func (this *ParameterInfo) GetTy() string {
	return this.Ty
}

func (this *ParameterInfo) GetDescription() string {
	return this.Description
}
