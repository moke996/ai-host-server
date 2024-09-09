package model

// http回复的消息内容
type HttpResponseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type EmptyRespData struct {
}

type StartRequest struct {
	Version string `json:"version"` // 版本号，为空表示获取所有版本号
	Tag     int    `json:"tag"`     // 每次对话的标识符
	Male    string `json:"content"`
	Female  string `json:"female"`
}

type RunRequest struct {
	Version  string `json:"version"` // 版本号，为空表示获取所有版本号
	Tag      int    `json:"tag"`     // 每次对话的标识符
	Content  string `json:"content"`
	NextStep string `json:"next"`
}

type SaveRequest struct {
	Title string `json:"title"`
	Info  string `json:"info"`
	Tag   int    `json:"tag"` // 每次对话的标识符
}

type HistoryRequest struct {
	Name string `json:"name"`
}

type HistoryResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"` // 版本号，为空表示获取所有版本号
	Tag     int    `json:"tag"`     // 每次对话的标识符
	Male    string `json:"content"` //用户profile
	Female  string `json:"female"`  // 用户profile
}
