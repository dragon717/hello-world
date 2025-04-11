package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var GmodelPool *ModelPool
var gctx context.Context

type ModelPool struct {
	ApiKeyList []string
	Pool       []*genai.GenerativeModel
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
		Pool:  make([]*genai.GenerativeModel, 0),
		Index: 0,
	}
	apikey := strings.Split(os.Getenv("APIKEY"), ";")
	for _, s := range apikey {
		pool.ApiKeyList = append(pool.ApiKeyList, s)
	}

	for r := 0; r < len(pool.ApiKeyList); r++ {
		client, err := genai.NewClient(gctx, option.WithAPIKey(pool.ApiKeyList[r]))
		if err != nil {
			log.Fatal(err)
		}
		model := client.GenerativeModel("gemini-2.0-flash-001")

		Temperature := float32(0.8)
		model.Temperature = &Temperature

		TopK := int32(40)
		model.TopK = &TopK

		pool.Pool = append(pool.Pool, model)
	}

	GmodelPool = pool
}

func sendmsg(entity EntityInterface) string {
	// 添加明确的JSON格式要求到prompt
	//m := fmt.Sprintf("你的信息:%s,ActionLog是你的记忆,地图信息(二维数组):%s,行为类型:%s,实体类型:%s,你是一个实体,你可以进行 行为类型 进行交互或者活动,你给出一个你想进行的行为,尽量不要重复,请用JSON格式响应以下请求，不要包含任何解释或格式化字符。JSON格式:{selfid:number,action:ACTION_TYPE,target:{x:number,y:number,id:number},reason:string},例子:{selfid:1001,action:100001,target:{x:0,y:2,id:8888},reason:'你对坐标[0,2]的实体id为8888的树进行砍树'}", common.JsonMarshal(entity), common.JsonMarshal(Gmap.Map), common.JsonMarshal(ActionParamCfg), common.JsonMarshal(GParamCfg))
	m := formatActionMsg(entity)
	if !*(devMode) {
		//fmt.Println(m)
	}
	prompt := genai.Text(m)

	resp, err := GmodelPool.Get().GenerateContent(gctx, prompt)
	if err != nil {
		if !*(devMode) {
			log.Println("[%d] GenerateContent error:", GmodelPool.Index, err)
		}
		sendmsg(entity)
		return ""
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
	entity.SendActionChan(formatRes(jsonData))

	js, _ := json.Marshal(jsonData)
	fmt.Printf("[%v]---[ACTION]---%v", time.Now(), string(js))
	return string(js)
}
func formatRes(actionmsg *ActionMsg) *ActionMsg {
	if actionmsg.Target.Item == nil {
		actionmsg.Target.Item = &ItemTarget{}
	}
	if actionmsg.Target.Position == nil {
		actionmsg.Target.Position = &Position{}
	}
	if actionmsg.Target.Entity == nil {
		actionmsg.Target.Entity = &EntityTarget{}
	}
	if actionmsg.Target.Craft == nil {
		actionmsg.Target.Craft = &CraftTarget{}
	}
	if actionmsg.Target.Duration == nil {
		actionmsg.Target.Duration = &DurationTarget{}
	}
	return actionmsg
}
