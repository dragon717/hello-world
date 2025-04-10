package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
)

var GmodelPool *ModelPool
var gctx context.Context
var devMode = flag.Bool("dev", false, "Enable development mode (disable printing)")

type ModelPool struct {
	ApiKeyList []string
	Pool       []*genai.GenerativeModel
	Index      int
}

func (m *ModelPool) Get() *genai.GenerativeModel {
	for i := 0; i < len(m.Pool); i++ {
		p := m.Pool[m.Index]
		m.Index++
		if m.Index >= len(m.Pool) {
			m.Index = 0
		}
		return p
	}
	return nil
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

	for _, s := range pool.ApiKeyList {
		client, err := genai.NewClient(gctx, option.WithAPIKey(s))
		if err != nil {
			log.Fatal(err)
		}
		model := client.GenerativeModel("gemini-2.0-flash-001")
		pool.Pool = append(pool.Pool, model)
	}

	GmodelPool = pool
}

func sendmsg(entity EntityInterface) string {
	// 添加明确的JSON格式要求到prompt
	//m := fmt.Sprintf("你的信息:%s,ActionLog是你的记忆,地图信息(二维数组):%s,行为类型:%s,实体类型:%s,你是一个实体,你可以进行 行为类型 进行交互或者活动,你给出一个你想进行的行为,尽量不要重复,请用JSON格式响应以下请求，不要包含任何解释或格式化字符。JSON格式:{selfid:number,action:ACTION_TYPE,target:{x:number,y:number,id:number},reason:string},例子:{selfid:1001,action:100001,target:{x:0,y:2,id:8888},reason:'你对坐标[0,2]的实体id为8888的树进行砍树'}", common.JsonMarshal(entity), common.JsonMarshal(Gmap.Map), common.JsonMarshal(ActionParamCfg), common.JsonMarshal(GParamCfg))
	m := formatActionMsg(entity)
	if !(*devMode) {
		fmt.Println(m)
	}
	prompt := genai.Text(m)

	var resp *genai.GenerateContentResponse
	var err error
	for i := 0; i < len(GmodelPool.Pool); i++ {
		model := GmodelPool.Get()
		if model == nil {
			log.Println("所有 API Key 均已达到速率限制，无法生成内容。")
			return ""
		}
		resp, err = model.GenerateContent(gctx, prompt)
		if err != nil {
			if strings.Contains(err.Error(), "Error 429: You exceeded your current quota") {
				log.Println("GenerateContent error:", err)
				continue
			} else {
				log.Fatal("GenerateContent error:", err)
			}
		} else {
			break
		}
	}

	if err != nil {
		log.Println("所有 API Key 尝试均失败，无法生成内容。")
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
	entity.SendActionChan(jsonData)
	target := FindTarget(uint32(jsonData.Target.Id))
	if target != nil && jsonData.Target.Id != jsonData.SelfId {
		target.SendActionChan(jsonData)
	}

	js, _ := json.Marshal(jsonData)

	return string(js)
}
