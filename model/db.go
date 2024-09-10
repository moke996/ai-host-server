package model

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	AiHostPrompt  = "ai_host_prompt"
	AiTestHistory = "ai_test_history"
	AiUserProfile = "ai_host_profile"
)

type Prompt struct {
	Id                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type               string             `json:"type,omitempty" bson:"type"`
	Version            string             `json:"version" bson:"version"`
	System             string             `json:"system"`
	Introduction       string             `json:"introduction" bson:"introduction"`
	IcebreakerQuestion string             `json:"icebreaker_question" bson:"icebreaker_question"`
	FollowUpQuestion   string             `json:"follow_up_question" bson:" follow_up_question"`
	NewQuestion        string             `json:"new_question" bson:"new_question"`
	Banter             string             `json:"banter" bson:"banter"`
	Reply              string             `json:"reply" bson:"reply"`
	WarnUp             string             `json:"warn_up" bson:"warn_up"`
	UserAnswer         string             `json:"user_answer" bson:"user_answer"`
	FinalRound         string             `json:"final_round" bson:"final_round"`
}

type History struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Version     string             `json:"version" bson:"version"`
	Male        string             `json:"male" bson:"male"`
	Female      string             `json:"female" bson:"female"`
	Tag         int                `json:"tag,omitempty" bson:"tag,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Content     string             `json:"content,,omitempty" bson:"content"`
	SaveTime    int64              `json:"save_time" bson:"save_time"`
}

type PromptConfig struct {
	Prompt   string
	Function string
}
