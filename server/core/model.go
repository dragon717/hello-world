package main

import (
	"Test/common"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
)

var GmodelPool *ModelPool
var gctx context.Context

type ModelPool struct {
	ApiKeyList []string
	Pool       []*genai.GenerativeModel
	Size       int
	Index      int
}

func (m *ModelPool) Get() *genai.GenerativeModel {
	p := m.Pool[m.Index]
	m.Index++
	if m.Index >= len(m.Pool) {
		m.Index = 0
	}
	return p
}
func initModelPool(size int) {
	gctx = context.Background()
	pool := &ModelPool{
		Pool:  make([]*genai.GenerativeModel, size),
		Size:  size,
		Index: 0,
	}
	apikey := strings.Split(os.Getenv("APIKEY"), ";")
	for _, s := range apikey {
		pool.ApiKeyList = append(pool.ApiKeyList, s)
	}

	for r := 0; r < len(pool.ApiKeyList); r++ {
		client, err := genai.NewClient(gctx, option.WithAPIKey("AIzaSyAGt3wFGsSE0AJALtItWnd1S_QKZ5cBvwo"))
		if err != nil {
			log.Fatal(err)
		}
		model := client.GenerativeModel("gemini-2.0-flash-001")
		pool.Pool[r] = model
	}

	GmodelPool = pool
}

func sendmsg(entity EntityInterface) string {
	// 添加明确的JSON格式要求到prompt
	//m := fmt.Sprintf("你的信息:%s,ActionLog是你的记忆,地图信息(二维数组):%s,行为类型:%s,实体类型:%s,你是一个实体,你可以进行 行为类型 进行交互或者活动,你给出一个你想进行的行为,尽量不要重复,请用JSON格式响应以下请求，不要包含任何解释或格式化字符。JSON格式:{selfid:number,action:ACTION_TYPE,target:{x:number,y:number,id:number},reason:string},例子:{selfid:1001,action:100001,target:{x:0,y:2,id:8888},reason:'你对坐标[0,2]的实体id为8888的树进行砍树'}", common.JsonMarshal(entity), common.JsonMarshal(Gmap.Map), common.JsonMarshal(ActionParamCfg), common.JsonMarshal(GParamCfg))
	m := formatActionMsg(entity)
	fmt.Println(m)
	prompt := genai.Text(m)

	resp, err := GmodelPool.Get().GenerateContent(gctx, prompt)
	if err != nil {
		log.Fatal("GenerateContent error:", err)
	}

	var rawResponse strings.Builder
	if resp.Candidates != nil {
		for _, candidate := range resp.Candidates {
			for _, part := range candidate.Content.Parts {
				rawResponse.WriteString(string(part.(genai.Text)))
			}
		}
	}

	// 验证并清理JSON格式
	var jsonData *ActionMsg
	s1 := strings.ReplaceAll(rawResponse.String(), "```", "")
	s2 := strings.ReplaceAll(s1, "json", "")
	if err := json.Unmarshal([]byte(s2), &jsonData); err != nil {
		log.Fatalf("响应格式无效: %v\n原始响应: %s", err, rawResponse.String())
	}
	entity.SendActionChan(jsonData)
	target := FindTarget(uint32(jsonData.Target.Id))
	if target != nil && jsonData.Target.Id != jsonData.SelfId {
		target.SendActionChan(jsonData)
	}

	js, _ := json.Marshal(jsonData)

	return string(js)
}

func formatActionMsg(you EntityInterface) string {
	var msg string

	var record string
	for _, actionLog := range you.GetActionLog() {
		record += fmt.Sprintf("%s %s %s,", actionLog.Time, actionLog.Action, actionLog.Result)
	}
	//record, _ := json.Marshal(you.GetActionLog())
	baseinfo := fmt.Sprintf("你的ID:%d;名字:%s;类型ID:%d;背包:%+v;坐标:[%d,%d],生命值:%d;饱食度:%d;记忆:(%s) ", you.GetId(), you.GetName(), you.GetType(), you.GetBag(), you.GetX(), you.GetY(), you.GetHP(), you.GetSatietyDegree(), string(record))
	mapInfo := fmt.Sprintf("当前时间:%s,地图周围信息:%s", WorldMap.Gmap.GlobalTime.GetTime(), WorldMap.Gmap.GetAroundInfo(int(you.GetX()), int(you.GetY()), 1))
	actionInfo := fmt.Sprintf("行为类型:%s", common.JsonMarshal(ActionParamCfg))
	entityInfo := fmt.Sprintf("实体类型:%s", common.JsonMarshal(GParamCfg))
	msg = fmt.Sprintf("%s,%s,%s,%s,你可以使用'行为类型'进行交互或者活动,你给出一个你想进行的行为,尽量不要重复,请用JSON格式响应以下请求，不要包含任何解释或格式化字符。JSON格式:{selfid:number,action:ACTION_TYPE,target:{x:number,y:number,id:number},reason:string},例子:{selfid:1001,action:100001,target:{x:0,y:2,id:8888},reason:''} 解释:你使用action:100001(ActionTypeCuttingDownTrees)对坐标[0,2]的id为8888的树进行了砍伐", baseinfo, mapInfo, actionInfo, entityInfo)
	return msg
}
