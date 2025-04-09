package main

import (
	"Test/common"
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"
)

type TemplateParameter struct {
	Template
}

func (this *TemplateParameter) GenGoParamFile(dstfile string, tarStruct interface{}) {
	dstfile = strings.Replace(dstfile, ".go", "_config.go", 1)
	if Exist(dstfile) {
		err := os.Remove(dstfile)
		if err != nil {
			fmt.Printf("Remove file err %s", err)
			return
		}
	}

	data := tarStruct.(*TarStruct)

	var pListdata = &ParamConfig{
		Name:       FormatName(data.OriginName),
		OriginName: data.OriginName,
		FileName:   dstfile,
		PList:      make([]*ParamList, 0),
	}
	err := common.LoadConfig(data.OriginPath, &pListdata)
	if err != nil {
		fmt.Println("xml.Unmarshal err", err)
		return
	}

	for _, list := range pListdata.PList {
		list.ValName = FormatName(list.Key)
		switch list.Type {
		case "uint32":
			list.Type = "uint32"
		case "string":
			list.Type = "string"
		case "uint32[]":
			list.Type = "[]uint32"
		case "[]uint32":
			list.Type = "[]uint32"
		case "string[]":
			list.Type = "[]string"
		case "[]string":
			list.Type = "[]string"
		default:
			panic("Failed to parse parameter type:" + list.Type)
		}
	}

	//创建ConfigMg.go文件
	file, err := os.OpenFile(dstfile, os.O_CREATE|os.O_WRONLY, os.FileMode(0666))
	if err != nil {
		fmt.Print("create config_manager.go fail, err:", err.Error())
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	err = this.ExecuteTemplate(buf, this.execTmpl, pListdata)
	if err != nil {
		fmt.Println("execute config parameter err", err)
		return
	}
	dst, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(dstfile, " format err", err)
		return
	}

	n, err := file.Write(dst)
	if err != nil {
		fmt.Println(dstfile, " write to file err", err)
		return
	}
	if n < len(dst) {
		fmt.Println(dstfile, " write to file less")
		return
	}
}
