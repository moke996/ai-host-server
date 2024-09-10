package controller

import (
	"ai-dating/constant"
	"ai-dating/model"
	"ai-dating/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
	"time"
)

var Prompt = make(map[string]map[string]*model.PromptConfig)

// 获取版本号
func GetVersion(c *gin.Context) {
	var result []string
	if len(Prompt) > 0 {
		for version, _ := range Prompt {
			result = append(result, version)
		}
		HttpSuccess(c, result)
		return
	}
	InitData()
	for version, _ := range Prompt {
		result = append(result, version)
	}
	HttpSuccess(c, result)
	return

}

// 开始
func Start(c *gin.Context) {
	var req model.StartRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		HttpFail(c, nil, "Start ShouldBindJSON error! "+err.Error())
		return
	}
	if req.Male == "" || req.Female == "" || req.Tag == 0 {
		HttpFail(c, nil, "Start params error! Male, Female and Tag must be provided!")
		return
	}

	prompt, err := getLatestPrompt(req.Version, constant.System)
	if err != nil {
		HttpFail(c, nil, "Start getLatestPrompt error! err: "+err.Error())
		return
	}

	// 替换profile
	item := strings.Replace(prompt.Prompt, "{{$male}}", req.Male, 1)
	prompt.Prompt = strings.Replace(item, "{{$female}}", req.Female, 1)

	msg := []model.OpenAIMessageType{
		{
			Role:    constant.System,
			Content: prompt.Prompt,
		},
		{
			Role:    constant.RoleUser,
			Content: "start",
		},
	}

	function := model.FunctionType{}
	err = json.Unmarshal([]byte(prompt.Function), &function)
	if err != nil {
		HttpFail(c, nil, "Start Unmarshal function error! err: "+err.Error())
		return
	}
	firstReply, err := OpenAIReply(c, msg, function, req.Tag)
	if err != nil {
		HttpFail(c, nil, "Start OpenAIReply error! err: "+err.Error())
		return
	}
	resp := &model.AIResponse{}
	if firstReply.NextStep == constant.Introduction {
		params := &model.RunRequest{
			Version:  req.Version,
			Tag:      req.Tag,
			Content:  "",
			NextStep: constant.Introduction,
		}
		resp, err = GetOpenAIReply(c, params)
		if err != nil {
			HttpFail(c, nil, "Run GetOpenAIReply error! err: "+err.Error())
			return
		}
	}
	HttpSuccess(c, resp)
	return
}

// Run
func Run(c *gin.Context) {
	var req *model.RunRequest
	err := c.ShouldBindJSON(req)
	if err != nil {
		HttpFail(c, nil, "Start ShouldBindJSON error! "+err.Error())
		return
	}
	if req.Tag == 0 || req.NextStep == "" {
		HttpFail(c, nil, "Start params error! next_step and Tag must be provided!")
		return
	}
	resp, err := GetOpenAIReply(c, req)
	if err != nil {
		HttpFail(c, nil, "Run GetOpenAIReply error! err: "+err.Error())
		return
	}
	HttpSuccess(c, resp)
	return

}

// 保存历史记录
func Save(c *gin.Context) {
	var req model.SaveRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		HttpFail(c, nil, "Start ShouldBindJSON error! "+err.Error())
		return
	}
	// 查询profile
	cache := repository.NewCacheHistory()
	profile, err := cache.GetUserProfile(c, req.Tag)
	if err != nil {
		HttpFail(c, nil, "GetUserProfile error! err: "+err.Error())
		return
	}
	strBytes, err := json.Marshal(req.Msg)
	if err != nil {
		HttpFail(c, nil, "Save Marshal history error! err: "+err.Error())
		return
	}
	err = repository.NewPrompt().InsertHistory(c, req, string(strBytes), profile.Male, profile.Female)
	if err != nil {
		HttpFail(c, nil, "Save InsertHistory error! err: "+err.Error())
		return
	}
	HttpSuccess(c, nil)
	return

}

// 获取历史记录列表
func GetHistoryList(c *gin.Context) {
	data, err := repository.NewPrompt().GetHistoryList(c)
	if err != nil {
		HttpFail(c, nil, "GetHistoryList error! err: "+err.Error())
		return
	}
	HttpSuccess(c, data)
	return
}

func GetHistoryById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		HttpFail(c, nil, "id is empty! ")
		return
	}
	idObj, err := Str2ObjectId(id)
	if err != nil {
		HttpFail(c, nil, " Str2ObjectId id is err! "+err.Error())
		return
	}

	history, err := repository.NewPrompt().GetHistoryById(c, idObj)
	if err != nil {
		HttpFail(c, nil, "GetHistoryById error! err: "+err.Error())
		return
	}

	msg := make([]model.MessageResp, 0)
	err = json.Unmarshal([]byte(history.Content), &msg)
	if err != nil {
		HttpFail(c, nil, "Unmarshal history.Content error! err: "+err.Error())
		return
	}
	resp := model.HistoryResponse{
		Messages: msg,
		Version:  history.Version,
		Profile: model.Profile{
			Male:   history.Male,
			Female: history.Female,
		},
	}
	HttpSuccess(c, resp)
	return
}

func OpenAIReply(c *gin.Context, msg []model.OpenAIMessageType, function model.FunctionType, tag int) (*model.AIResponse, error) {
	for i := 0; i < 10; i++ {
		data, err := OpenAiReply(msg, function)
		if err != nil {
			time.Sleep(50 * time.Millisecond) // 等待一段时间后重试
			continue
		}
		if data.NextStep == "" {
			fmt.Println(" OpenAIReply nextstep is empty!: ")
			continue
		}
		if data.A != "" && data.B != "" {
			data.AiAnswer = data.AiAnswer + "\n" + data.A + "\n" + data.B
			data.A = ""
			data.B = ""
		}
		// AI回复加入到会话历史
		strByte, _ := json.Marshal(data)
		msg = append(msg, model.OpenAIMessageType{
			Role:    constant.RoleAssistant,
			Content: string(strByte),
		})
		r := repository.NewCacheHistory()
		err = r.AddSessionHistory(c, strconv.Itoa(tag), msg)
		if err != nil {
			fmt.Println(" AddSessionHistory error! err: ", err)
		}
		return data, nil
	}
	return nil, errors.New("openAIReply no answer! ")
}

// 初始化数据
func InitData() {
	ctx := context.Background()
	prompt := repository.NewPrompt()
	promptMap, err := prompt.GetPromptList(ctx, constant.AiPrompt)
	if err != nil {
		fmt.Println(" InitData GetPromptList error! err: ", err)
		return
	}
	functionMap, err := prompt.GetPromptList(ctx, constant.AiFunctions)
	if err != nil {
		fmt.Println(" InitData GetFunction error! err: ", err)
		return
	}
	for version, val := range promptMap {
		params := make(map[string]*model.PromptConfig)
		functions, ok := functionMap[version]
		if !ok {
			functions = functionMap["v1.0"]
		}
		params[constant.System] = &model.PromptConfig{
			Prompt:   val.System,
			Function: functions.System,
		}
		params[constant.Introduction] = &model.PromptConfig{
			Prompt:   val.Introduction,
			Function: functions.Introduction,
		}
		params[constant.IcebreakerQuestion] = &model.PromptConfig{
			Prompt:   val.IcebreakerQuestion,
			Function: functions.IcebreakerQuestion,
		}
		params[constant.FollowUpQuestion] = &model.PromptConfig{
			Prompt:   val.FollowUpQuestion,
			Function: functions.FollowUpQuestion,
		}
		params[constant.NewQuestion] = &model.PromptConfig{
			Prompt:   val.NewQuestion,
			Function: functions.NewQuestion,
		}
		params[constant.Banter] = &model.PromptConfig{
			Prompt:   val.Banter,
			Function: functions.Banter,
		}
		params[constant.Reply] = &model.PromptConfig{
			Prompt:   val.Reply,
			Function: functions.Reply,
		}
		params[constant.WarnUp] = &model.PromptConfig{
			Prompt:   val.WarnUp,
			Function: functions.WarnUp,
		}
		params[constant.UserAnswer] = &model.PromptConfig{
			Prompt:   val.UserAnswer,
			Function: functions.UserAnswer,
		}
		params[constant.FinalRound] = &model.PromptConfig{
			Prompt:   val.FinalRound,
			Function: functions.FinalRound,
		}
		fmt.Println(Prompt)
		Prompt[version] = params
	}
}

func getLatestPrompt(version, key string) (*model.PromptConfig, error) {
	if prompt, ok := Prompt[version]; !ok {
		return nil, errors.New("version not found!" + version)
	} else {
		if v, ok := prompt[key]; !ok {
			return nil, errors.New(" prompt key not found!" + version)
		} else {
			return v, nil
		}
	}
}

func DeleteHistoryById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		HttpFail(c, nil, "id is empty! ")
		return
	}
	idObj, err := Str2ObjectId(id)
	if err != nil {
		HttpFail(c, nil, " Str2ObjectId id is err! "+err.Error())
		return
	}

	err = repository.NewPrompt().DeleteHistory(c, idObj)
	if err != nil {
		HttpFail(c, nil, "GetHistoryById error! err: "+err.Error())
		return
	}
	HttpSuccess(c, "")
	return
}

func Str2ObjectId(id string) (objId primitive.ObjectID, err error) {

	if id == "" {
		return primitive.NewObjectID(), errors.New("id is empty")
	}

	if objId, err = primitive.ObjectIDFromHex(id); err != nil {
		return objId, errors.New(err.Error())
	}
	return
}

func GetOpenAIReply(c *gin.Context, req *model.RunRequest) (*model.AIResponse, error) {
	cache := repository.NewCacheHistory()
	// 查询本轮使用的prompt
	prompt, err := getLatestPrompt(req.Version, req.NextStep)
	if err != nil {
		return nil, errors.New("Run getLatestPrompt error! err: " + err.Error())
	}
	// 查询会话历史
	history, err := cache.GetSessionHistory(c, strconv.Itoa(req.Tag))
	if err != nil {
		return nil, errors.New("Run GetSessionHistory error! err: " + err.Error())
	}
	if req.NextStep == constant.UserAnswer {
		history = append(history, model.OpenAIMessageType{
			Role:    constant.RoleUser,
			Content: req.Content,
		})
	} else {
		history = append(history, model.OpenAIMessageType{
			Role:    constant.RoleUser,
			Content: prompt.Prompt,
		})
	}

	function := model.FunctionType{}
	err = json.Unmarshal([]byte(prompt.Function), &function)
	if err != nil {
		return nil, errors.New("Start Unmarshal function error! err: " + err.Error())
	}
	resp, err := OpenAIReply(c, history, function, req.Tag)
	if err != nil {
		return nil, errors.New("Start OpenAIReply error! err: " + err.Error())
	}
	return resp, nil
}
