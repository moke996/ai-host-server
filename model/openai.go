package model

type OpenAIRequest struct {
	Model        string              `json:"model"`
	Messages     []OpenAIMessageType `json:"messages"`
	Functions    []FunctionType      `json:"functions"`
	FunctionCall map[string]string   `json:"function_call,omitempty"`
}

type OpenAIMessageType struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type FunctionType struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Required    []string               `json:"required"`
	Type        string                 `json:"type"`
}

type OpenAIResponse struct {
	Error   interface{} `json:"error,omitempty"`
	Choices []struct {
		Message struct {
			Role         string `json:"role"`
			Content      string `json:"content"`
			FunctionCall struct {
				Name      string      `json:"name"`
				Arguments interface{} `json:"arguments"`
			} `json:"function_call,omitempty"`
		} `json:"message"`
	} `json:"choices"`
}
