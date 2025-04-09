package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"text/template"
)

type Template struct {
	*template.Template
	tmplDir  string
	execTmpl string
}

func (this *Template) Init(tmplDir, execTmpl string) {
	this.tmplDir = tmplDir
	this.execTmpl = execTmpl
	var name = tmplDir + execTmpl
	this.Template = template.Must(template.ParseFiles(name))
}

// func (this *Template) GenGoFile(templ_file, dstfile string, tarStruct interface{}) {
func (this *Template) GenGoFile(dstfile string, tarStruct interface{}) {
	if Exist(dstfile) {
		err := os.Remove(dstfile)
		if err != nil {
			//如果删除失败则输出 file remove Error!
			fmt.Println("file remove Error!")
			//输出错误详细信息
			fmt.Printf("%s", err)
			return
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
	err = this.ExecuteTemplate(buf, this.execTmpl, tarStruct)
	if err != nil {
		fmt.Println("execute template err", err)
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
