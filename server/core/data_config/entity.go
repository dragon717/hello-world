// 该文件为工具生成 不要修改
package data_config

import (
	"Test/common"
	"unsafe"
)

func GetXmlEntityInfo() *XmlEntityInfo {
	data_sheet, err := LoadXmlEntityInfo()
	if err != nil {
		return nil
	}
	return data_sheet
}

type XmlEntityInfo struct {
	Datas []*EntityInfo
}

func GetXmlEntityInfoName() string {
	return DataConfigDir + "entity.xml"
}
func LoadXmlEntityInfo() (*XmlEntityInfo, error) {
	return LoadXmlEntityInfoByFileName(DataConfigDir + "entity.xml")
}

func LoadXmlEntityInfoByFileName(filename string) (*XmlEntityInfo, error) {
	data_sheet := &struct {
		Datas []*struct {
			Name        string `xml:"name,attr"`
			Id          int    `xml:"id,attr"`
			TypeId      int    `xml:"type_id,attr"`
			Description string `xml:"description,attr"`
		} `xml:"data"`
	}{}

	if err := common.LoadConfig(filename, data_sheet); err != nil {
		return nil, err
	}
	return (*XmlEntityInfo)(unsafe.Pointer(data_sheet)), nil
}

type EntityInfo struct {
	Name        string // 序号
	Id          int    // 参数内容
	TypeId      int    // 参数类型
	Description string // 细节描述
}

func (this *EntityInfo) GetName() string {
	return this.Name
}

func (this *EntityInfo) GetId() int {
	return this.Id
}

func (this *EntityInfo) GetTypeId() int {
	return this.TypeId
}

func (this *EntityInfo) GetDescription() string {
	return this.Description
}
