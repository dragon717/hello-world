package common

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
)

func LoadConfig(filename string, v interface{}) error {
	contents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		if err = xml.Unmarshal(contents, v); err != nil {
			return err
		}
		return nil
	}
}
func JsonMarshal(s interface{}) string {
	marshal, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json Marshal err", err)
		return ""
	}
	return string(marshal)
}
