package main

type ParamList struct {
	ValName     string
	Key         string `xml:"key,attr"`
	Content     string `xml:"content,attr"`
	Type        string `xml:"ty,attr"`
	Description string `xml:"description,attr"`
}

type ParamConfig struct {
	Name       string
	OriginName string
	FileName   string
	PList      []*ParamList `xml:"data"`
}
