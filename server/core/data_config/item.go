// 该文件为工具生成 不要修改
package data_config

import (
	"Test/common"
	"unsafe"
)

func GetXmlItemInfo() *XmlItemInfo {
	data_sheet, err := LoadXmlItemInfo()
	if err != nil {
		return nil
	}
	return data_sheet
}

type XmlItemInfo struct {
	Datas []*ItemInfo
}

func GetXmlItemInfoName() string {
	return DataConfigDir + "item.xml"
}
func LoadXmlItemInfo() (*XmlItemInfo, error) {
	return LoadXmlItemInfoByFileName(DataConfigDir + "item.xml")
}

func LoadXmlItemInfoByFileName(filename string) (*XmlItemInfo, error) {
	data_sheet := &struct {
		Datas []*struct {
			Name        string `xml:"name,attr"`
			Id          int    `xml:"id,attr"`
			TypeId      int    `xml:"type_id,attr"`
			Description string `xml:"description,attr"`
			NeedKey1    int    `xml:"need_key1,attr"`
			NeedValue1  int    `xml:"need_value1,attr"`
			NeedKey2    int    `xml:"need_key2,attr"`
			NeedValue2  int    `xml:"need_value2,attr"`
			NeedKey3    int    `xml:"need_key3,attr"`
			NeedValue3  int    `xml:"need_value3,attr"`
		} `xml:"data"`
	}{}

	if err := common.LoadConfig(filename, data_sheet); err != nil {
		return nil, err
	}
	return (*XmlItemInfo)(unsafe.Pointer(data_sheet)), nil
}

type ItemInfo struct {
	Name        string // 序号
	Id          int    // 参数id
	TypeId      int    // 参数类型
	Description string // 细节描述
	NeedKey1    int    // 合成需要的材料类型1
	NeedValue1  int    // 合成需要的材料数量1
	NeedKey2    int    // 合成需要的材料类型2
	NeedValue2  int    // 合成需要的材料数量2
	NeedKey3    int    // 合成需要的材料类型3
	NeedValue3  int    // 合成需要的材料数量3
}

func (this *ItemInfo) GetName() string {
	return this.Name
}

func (this *ItemInfo) GetId() int {
	return this.Id
}

func (this *ItemInfo) GetTypeId() int {
	return this.TypeId
}

func (this *ItemInfo) GetDescription() string {
	return this.Description
}

func (this *ItemInfo) GetNeedKey1() int {
	return this.NeedKey1
}

func (this *ItemInfo) GetNeedValue1() int {
	return this.NeedValue1
}

func (this *ItemInfo) GetNeedKey2() int {
	return this.NeedKey2
}

func (this *ItemInfo) GetNeedValue2() int {
	return this.NeedValue2
}

func (this *ItemInfo) GetNeedKey3() int {
	return this.NeedKey3
}

func (this *ItemInfo) GetNeedValue3() int {
	return this.NeedValue3
}
