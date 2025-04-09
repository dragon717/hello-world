package main

type ActionMsg struct {
	SelfId int32   `json:"selfid"`
	Action uint32  `json:"action"`
	Target *Target `json:"target"`
	Reason string  `json:"reason"`
}
type Target struct {
	X  int   `json:"x"`
	Y  int   `json:"y"`
	Id int32 `json:"id"`
}

type ResultMsg struct {
	ActionID  uint32 `json:"action"`
	Result    string `json:"result"`
	AwardList map[uint32]uint32
}
