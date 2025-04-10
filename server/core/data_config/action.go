// 该文件为工具生成 不要修改
package data_config

import (
	"Test/common"
	"unsafe"
)

func GetXmlActionInfo() *XmlActionInfo {
	data_sheet, err := LoadXmlActionInfo()
	if err != nil {
		return nil
	}
	return data_sheet
}

type XmlActionInfo struct {
	Datas []*ActionInfo
}

func GetXmlActionInfoName() string {
	return DataConfigDir + "action.xml"
}
func LoadXmlActionInfo() (*XmlActionInfo, error) {
	return LoadXmlActionInfoByFileName(DataConfigDir + "action.xml")
}

func LoadXmlActionInfoByFileName(filename string) (*XmlActionInfo, error) {
	data_sheet := &struct {
		Datas []*struct {
			Name        string `xml:"name,attr"`
			Id          int    `xml:"id,attr"`
			Ty          string `xml:"ty,attr"`
			TypeId      int    `xml:"type_id,attr"`
			Description string `xml:"description,attr"`
			NeedType    int    `xml:"need_type,attr"`
			NeedId      int    `xml:"need_id,attr"`
		} `xml:"data"`
	}{}

	if err := common.LoadConfig(filename, data_sheet); err != nil {
		return nil, err
	}
	return (*XmlActionInfo)(unsafe.Pointer(data_sheet)), nil
}

type ActionInfo struct {
	Name        string // 序号
	Id          int    // 参数内容
	Ty          string // 参数类型
	TypeId      int    // 参数类型
	Description string // 细节描述
	NeedType    int    // 执行需要的道具类型
	NeedId      int    // 执行需要的道具ID
}

func (this *ActionInfo) GetName() string {
	return this.Name
}

func (this *ActionInfo) GetId() int {
	return this.Id
}

func (this *ActionInfo) GetTy() string {
	return this.Ty
}

func (this *ActionInfo) GetTypeId() int {
	return this.TypeId
}

func (this *ActionInfo) GetDescription() string {
	return this.Description
}

func (this *ActionInfo) GetNeedType() int {
	return this.NeedType
}

func (this *ActionInfo) GetNeedId() int {
	return this.NeedId
}
