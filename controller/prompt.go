package controller

import (
	"ai-host/constant"
	"ai-host/model"
	"ai-host/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
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
	prompt.Prompt = strings.Replace(item, "{{$female}}", req.Male, 1)

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
	resp, err := OpenAIReply(c, msg, function, req.Tag)
	if err != nil {
		HttpFail(c, nil, "Start OpenAIReply error! err: "+err.Error())
		return
	}
	HttpSuccess(c, resp)
	return
}

// 初始化数据
func InitData() {
	ctx := context.Background()
	prompt := repository.NewPrompt()
	promptMap, err := prompt.GetPromptList(ctx)
	if err != nil {
		fmt.Println(" InitData GetPromptList error! err: ", err)
		return
	}
	functions, err := prompt.GetFunction(ctx)
	if err != nil {
		fmt.Println(" InitData GetFunction error! err: ", err)
		return
	}
	for version, val := range promptMap {
		params := make(map[string]*model.PromptConfig)
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
		params[constant.AfterAnswer] = &model.PromptConfig{
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

// Run
func Run(c *gin.Context) {
	var req model.RunRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		HttpFail(c, nil, "Start ShouldBindJSON error! "+err.Error())
		return
	}
	if req.Tag == 0 || req.NextStep == "" {
		HttpFail(c, nil, "Start params error! next_step and Tag must be provided!")
		return
	}
	cache := repository.NewCacheHistory()
	// 查询本轮使用的prompt
	prompt, err := getLatestPrompt(req.Version, req.NextStep)
	if err != nil {
		HttpFail(c, nil, "Run getLatestPrompt error! err: "+err.Error())
		return
	}
	// 查询会话历史
	history, err := cache.GetSessionHistory(c, strconv.Itoa(req.Tag))
	if err != nil {
		HttpFail(c, nil, "Run GetSessionHistory error! err: "+err.Error())
		return
	}
	if req.NextStep != constant.AfterAnswer {
		history = append(history, model.OpenAIMessageType{
			Role:    constant.RoleUser,
			Content: prompt.Prompt,
		})
	} else {
		history = append(history, model.OpenAIMessageType{
			Role:    constant.RoleUser,
			Content: req.Content,
		})
	}

	function := model.FunctionType{}
	err = json.Unmarshal([]byte(prompt.Function), &function)
	if err != nil {
		HttpFail(c, nil, "Start Unmarshal function error! err: "+err.Error())
		return
	}
	resp, err := OpenAIReply(c, history, function, req.Tag)
	if err != nil {
		HttpFail(c, nil, "Start OpenAIReply error! err: "+err.Error())
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
	// 查询会话历史
	cache := repository.NewCacheHistory()
	history, err := cache.GetSessionHistory(c, strconv.Itoa(req.Tag))
	if err != nil {
		HttpFail(c, nil, "Run GetSessionHistory error! err: "+err.Error())
		return
	}
	if len(history) == 0 {
		HttpFail(c, nil, "history not fount!")
		return
	}

	strBytes, err := json.Marshal(history)
	if err != nil {
		HttpFail(c, nil, "Save Marshal history error! err: "+err.Error())
		return
	}

	err = repository.NewPrompt().InsertHistory(c, req, string(strBytes))
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
	HttpSuccess(c, history)
	return
}

func OpenAIReply(c *gin.Context, msg []model.OpenAIMessageType, function model.FunctionType, tag int) (map[string]string, error) {
	for i := 0; i < 5; i++ {
		result, err := OpenAiWithRetry(msg, function, 3)
		if err != nil {
			fmt.Println(" OpenAIReply error! err: " + err.Error())
			continue
		}
		if _, ok := result[constant.NextStep]; !ok {
			fmt.Println(" OpenAIReply not support!")
			continue
		} else {
			strByte, _ := json.Marshal(result)
			msg = append(msg, model.OpenAIMessageType{
				Role:    constant.RoleAssistant,
				Content: string(strByte),
			})
			r := repository.NewCacheHistory()
			err := r.AddSessionHistory(c, strconv.Itoa(tag), msg)
			if err != nil {
				fmt.Println(" AddSessionHistory error! err: ", err)
			}
			return result, nil
		}
	}
	return nil, errors.New("openAIReply no answer! ")
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
