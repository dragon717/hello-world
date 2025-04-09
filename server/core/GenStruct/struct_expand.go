package main

// 扩展结构体字段
func ExpandStruct(src, target *TarStruct) {
	// 单向扩充-新字段追加到后面
	for idx, fieldName := range src.FieldsName {
		if targetIdx := target.GetFieldsIdxByName(fieldName); targetIdx >= 0 {
			continue
		}
		target.InsertFields(src, idx)
	}
	// 调整顺序-前面的字段名一样
	// 后面的字段名一样-
	// 没有字段名一样放最后
	//
}

func (this *TarStruct) GetFieldsIdxByName(name string) int {
	for idx, f := range this.FieldsName {
		if f == name {
			return idx
		}
	}
	return -1
}
func (this *TarStruct) InsertFields(from *TarStruct, idx int) {
	this.FieldsName = append(this.FieldsName, from.FieldsName[idx])
	this.FieldsType = append(this.FieldsType, from.FieldsType[idx])
	this.FieldsTag = append(this.FieldsTag, from.FieldsTag[idx])
	this.FieldsNote = append(this.FieldsNote, from.FieldsNote[idx]+" 扩充字段用不到")
}

func ExpandStructBoth(list []*TarStruct, name1, name2 string) {
	var src = FindTarStructByName(list, name1)
	var target = FindTarStructByName(list, name2)
	if src == nil {
		return
	}
	if target == nil {
		return
	}
	ExpandStruct(src, target)
	ExpandStruct(target, src)

	// 调整顺序
	for i := 0; i < len(src.FieldsName); i++ {
		if src.FieldsName[i] == target.FieldsName[i] {
			continue
		}
		// 找个一样的
		target.GetFieldsIdxByName(src.FieldsName[i])
		src.GetFieldsIdxByName(target.FieldsName[i])
	}
}

// 结构体
func FindTarStructByName(list []*TarStruct, name string) *TarStruct {
	for _, val := range list {
		if val.OriginName == name {
			return val
		}
	}
	return nil
}
