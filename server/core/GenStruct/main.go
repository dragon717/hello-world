package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type TarStruct struct {
	OriginPath string   //原始路径
	OriginName string   //原始名字
	Name       string   //结构体名称
	FieldsName []string //字段名称
	FieldsNote []string //字段注释
	FieldsType []string //字段类型
	FieldsTag  []string //字段标记
}

// 获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) ([]string, error) {
	files := make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 进入下一目录
			continue
		}
		if IsIgnore(fi.Name()) {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

func ReadFiles(files []string) []*TarStruct {
	var tarStruct = make([]*TarStruct, 0, len(files))
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return tarStruct
		}
		tstruct := &TarStruct{}
		tstruct.OriginPath = file
		tstruct.OriginName = GetOriginFileName(file)
		tstruct.Name = FormatName(tstruct.OriginName + "_info")
		tstruct.FieldsName = make([]string, 0)
		tstruct.FieldsType = make([]string, 0)
		tstruct.FieldsTag = make([]string, 0)
		buf := bufio.NewReader(f)
		stop := false
		con := false
		for !stop {
			line, err := buf.ReadString('\n')
			if err != nil {
				fmt.Println("读文件出错,file:", tstruct.Name)
				con = true
				break
			}
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "<!--") {
				stop = true
				strs := strings.Split(line, " ")
				for _, str := range strs {
					if str == "<!--" || str == "-->" || str == "" {
						continue
					}
					fieldsInfos := strings.Split(str, "=")
					if len(fieldsInfos) == 3 {
						tstruct.FieldsName = append(tstruct.FieldsName, FormatName(fieldsInfos[0]))
						tstruct.FieldsTag = append(tstruct.FieldsTag, fieldsInfos[0])
						tstruct.FieldsNote = append(tstruct.FieldsNote, fieldsInfos[1])
						switch fieldsInfos[2] {
						case "uint":
							tstruct.FieldsType = append(tstruct.FieldsType, "uint32")
						case "uint64":
							tstruct.FieldsType = append(tstruct.FieldsType, "uint64")
						case "int":
							tstruct.FieldsType = append(tstruct.FieldsType, "int")
						case "int64":
							tstruct.FieldsType = append(tstruct.FieldsType, "uint64")
						case "string":
							tstruct.FieldsType = append(tstruct.FieldsType, "string")
						case "int[]":
							tstruct.FieldsType = append(tstruct.FieldsType, "string")
						case "uint[]":
							tstruct.FieldsType = append(tstruct.FieldsType, "string")
						case "float":
							tstruct.FieldsType = append(tstruct.FieldsType, "float32")
						default:
							fmt.Println("不支持的类型,file:", tstruct.Name, " invalid type:", fieldsInfos[2])
							con = true
						}
					} else {
						fmt.Println("字段出错,file:", tstruct.Name, str, fieldsInfos)
						con = true
						break
					}
				}
			}
		}
		if con {
			continue
		}
		tarStruct = append(tarStruct, tstruct)
	}
	return tarStruct
}

func GetOriginFileName(path string) string {
	fileNames := strings.Split(filepath.Base(path), ".")
	if len(fileNames) == 2 {
		return fileNames[0]
	}
	return ""
}

func FormatName(filename string) string {
	var fileName string
	names := strings.Split(filename, "_")
	for _, name := range names {
		name = strings.ToLower(name)
		temp := make([]rune, 0)
		for index, ler := range []byte(name) {
			if index == 0 {
				if int(ler) >= 97 && int(ler) <= 122 {
					ler -= 32
				}
			}
			temp = append(temp, rune(ler))
		}
		fileName += string(temp)
	}
	return fileName
}

func FirstCharToLower(str string) string {
	if len(str) == 0 {
		return ""
	}

	b := []byte(str)
	if int(b[0]) >= 65 && int(b[0]) <= 90 {
		b[0] += 32
	}
	return string(b)
}

func CopyInitFile() {
	if !Exist(*InitGo) {
		fmt.Println("file: ", *InitGo, "not exist")
		return
	}

	b, err := ioutil.ReadFile(*InitGo)
	if err != nil {
		println("read file ", *InitGo, "err: ", err)
		return
	}

	f, err := os.OpenFile(*OutPath+"init.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		println("open file ", *OutPath+"init.go", "err: ", err)
		return
	}

	n, err := f.Write(b)
	if err != nil {
		println("write file ", *OutPath+"init.go", "err: ", err)
		return
	}
	if n < len(b) {
		println("write file ", *OutPath+"init.go", "less")
		return
	}
}

func GetAllFields(ts []*TarStruct) []string {
	fm := make(map[string]bool)
	res := make([]string, 0, 1000)

	for _, st := range ts {
		for _, fn := range st.FieldsName {
			if _, ok := fm[fn]; !ok {
				res = append(res, fn)
				fm[fn] = true
			}
		}
	}

	return res
}

func main() {
	fmt.Println("查找文件...")
	files, err := ListDir(*InPath, "xml")
	if err != nil {
		fmt.Println("读取文件失败")
		return
	}
	fmt.Printf("总共获取的文件:%d\n", len(files))
	fmt.Println("分析xml文件...")
	structs := ReadFiles(files)
	fmt.Println("错误文件数:", len(files)-len(structs))
	fmt.Println("生成go文件...")

	for _, v := range structs {
		out := strings.Join([]string{*OutPath, v.OriginName, ".go"}, "")
		if strings.HasSuffix(v.OriginName, "_parameter") {
			g_parameter_templ.GenGoFile(out, v)
		} else {
			g_templ.GenGoFile(out, v)
		}
	}

	for _, v := range structs {
		out := strings.Join([]string{*OutPath, v.OriginName, ".go"}, "")
		if strings.HasSuffix(v.OriginName, "parameter") || strings.HasSuffix(v.OriginName, "param") {
			g_parameter_static_templ.GenGoParamFile(out, v)
		}
	}

	fmt.Println("(: JOB DONE :)")
}
