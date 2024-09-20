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
	Male    string `json:"male"`
	Female  string `json:"female"`
}

type Profile struct {
	Male   string `json:"male"`
	Female string `json:"female"`
}

type RunRequest struct {
	Version  string `json:"version"` // 版本号，为空表示获取所有版本号
	Tag      int    `json:"tag"`     // 每次对话的标识符
	Content  string `json:"content"`
	NextStep string `json:"next"`
}

type AIResponse struct {
	NextStep string `json:"nextstep"`
	AiAnswer string `json:"ai_answer"`
	A        string `json:"A"`
	B        string `json:"B"`
}

type SaveRequest struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Version     string        `json:"version"`
	Tag         int           `json:"tag"` // 每次对话的标识符
	Message     []MessageResp `json:"message"`
}

type MessageResp struct {
	Text   string `json:"text"`
	Sender string `json:"sender"`
}

type HistoryResponse struct {
	Messages []MessageResp `json:"messages"`
	Version  string        `json:"version"`
	Profile
}

type GetPromptReq struct {
	Version string `json:"version"`
	Type    string `json:"type"` // 1: prompt  2: functions
	Tag     int    `json:"tag"`  // 每次对话的标识符
}

type ChildrenOptions struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type UpdatePromptReq struct {
	Version string `json:"version"`
	Type    string `json:"type"` // 1: prompt  2: functions
	Tag     string `json:"tag"`  // 每次对话的标识符
	Value   string `json:"value"`
}
