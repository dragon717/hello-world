package main

import (
	"bufio"
	"flag"
	"os"
	"path/filepath"
	"strings"
)

var (
	InPath                  = flag.String("in", "../data/", "数据配置原始目录")
	OutPath                 = flag.String("out", "../data_config/", "数据配置输出目录")
	ExecTempl               = flag.String("tmpl", "data_config_template.tmpl", "生成模板")
	ExecParameterTmpl       = flag.String("param_tmpl", "data_config_parameter.tmpl", "动态参数模板")
	ExecParameterStaticTmpl = flag.String("param_static_tmpl", "data_config_parameter_static.tmpl", "动态参数静态化模板")
	InitGo                  = flag.String("init", "./gen_templ/init.go", "数据配置init.go")
	TemplDir                = flag.String("tdir", "./gen_templ/", "模板目录")
	Ignore                  = flag.String("ig", "./xmlignore.txt", "忽略xml文件")
)

var (
	g_templ                  = &Template{}
	g_parameter_templ        = &Template{}
	g_parameter_static_templ = &TemplateParameter{}
	ignoreFiles              map[string]struct{}
	NullStruct               = struct{}{}
)

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func Mkdir(dir string) {
	pdir := filepath.Dir(dir)
	if !Exist(pdir) {
		Mkdir(pdir)
	}
	os.Mkdir(dir, os.FileMode(0777))
}

func InitIgnore() {
	ignoreFiles = make(map[string]struct{})
	if !Exist(*Ignore) {
		println("file ", *Ignore, "not exist")
		return
	}

	f, err := os.Open(*Ignore)
	if err != nil {
		println("open file ", *Ignore, "failed: ", err)
		return
	}

	rd := bufio.NewReader(f)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			break
		}
		lc := string(line)
		lc = strings.TrimSpace(lc)
		if lc != "" {
			println("ignore file ", lc)
			ignoreFiles[lc] = NullStruct
		}
	}
}

func IsIgnore(file string) bool {
	_, ok := ignoreFiles[file]
	return ok
}

func init() {
	flag.Parse()

	*TemplDir = strings.TrimSuffix(*TemplDir, "/") + "/"
	*OutPath = strings.TrimSuffix(*OutPath, "/") + "/"
	*InPath = strings.TrimSuffix(*InPath, "/") + "/"

	if !Exist(*TemplDir + *ExecTempl) {
		panic("template is not exist" + *TemplDir + *ExecTempl)
	}
	g_templ.Init(*TemplDir, *ExecTempl)
	g_parameter_templ.Init(*TemplDir, *ExecParameterTmpl)
	g_parameter_static_templ.Init(*TemplDir, *ExecParameterStaticTmpl)
	if !Exist(*OutPath) {
		Mkdir(*OutPath)
	}

}
