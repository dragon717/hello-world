/*
	该包提供数据配置读取。

使用者不需要实现从配置文件中读取数据。
使用者直接调用GetXml...即可获取数据配置。
默认数据配置读取目录为 ../data。
调用SetDataConfigDir(dir string) 设置数据配置目录。
*/
package data_config

import (
	"strings"
)

var (
	DataConfigDir = "../data/"
)

func SetDataConfigDir(dir string) {
	DataConfigDir = strings.TrimSuffix(dir, "/") + "/"
}
