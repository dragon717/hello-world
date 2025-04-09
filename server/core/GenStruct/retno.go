package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// 将retno.proto 中的错误码转成retno 包中的错误码

func GenRetno(fileName string, srcFileName string) error {
	var file, err = os.Open(fileName)
	if err != nil {
		return err
	}
	var content, err2 = io.ReadAll(file)
	if err2 != nil {
		return err2
	}
	var lines = strings.Split(string(content), "\n")
	var results []string
	for _, line := range lines {
		if strings.HasPrefix(line, "\tRET_") {
			// 增加uint32
			line = strings.Replace(line, "=", "uint32 = ", -1)
			// 去掉;
			line = strings.Replace(line, ";", "", -1)
			results = append(results, line)
		}
	}
	var str = fmt.Sprintf(`package retno
		const ( 
%s
)`, strings.Join(results, "\n"))
	return os.WriteFile(srcFileName, []byte(str), 0644)
}
