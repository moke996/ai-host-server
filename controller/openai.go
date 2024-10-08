package controller

import (
	"ai-dating/global"
	"ai-dating/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func OpenAiReply(msg []model.OpenAIMessageType, function model.FunctionType) (*model.AIResponse, error) {
	requestBody := model.OpenAIRequest{
		Model:        "gpt-4o",
		Messages:     msg,
		Functions:    []model.FunctionType{function},
		FunctionCall: map[string]string{"name": function.Name},
	}

	// 将请求体序列化为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("请求体序列化错误:", err)
		return nil, errors.New(" OpenAIReply json.Marshal failed! " + err.Error())
	}
	// 创建HTTP请求
	req, err := http.NewRequest("POST", global.Config.OpenAi.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("创建HTTP请求错误:", err)
		return nil, errors.New("OpenAIReply NewRequest failed! " + err.Error())
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+global.Config.OpenAi.Secret)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送HTTP请求错误:", err)
		return nil, errors.New("OpenAIReply do Request failed! " + err.Error())
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("OpenAIReply ReadAll failed! " + err.Error())
	}

	// 解析JSON响应
	var response model.OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New("OpenAIReply json.Unmarshal failed! body:" + string(body))
	}
	if response.Error != nil {
		d, _ := json.Marshal(response.Error)
		return nil, errors.New("OpenAIReply response error! body:" + string(d))
	}
	if response.Choices[0].Message.FunctionCall.Arguments == nil {
		ss, _ := json.Marshal(response)
		fmt.Println("Ai-Resp: ", string(ss))
		return nil, errors.New("OpenAIReply assistant is nil! ")
	}
	fmt.Println("AI-Req: ", function.Name)
	ss, _ := json.Marshal(response)
	fmt.Println("Ai-Resp: ", string(ss))
	assistant := response.Choices[0].Message.FunctionCall.Arguments
	result := &model.AIResponse{}
	err = json.Unmarshal([]byte(assistant.(string)), result)
	if err != nil {
		return nil, errors.New("OpenAIReply json.Unmarshal assistant failed! body:" + assistant.(string))
	}
	return result, nil
}
