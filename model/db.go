package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Prompt struct {
	Id                 primitive.ObjectID `json:"id" bson:"_id"`
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
	Id       primitive.ObjectID `json:"id" bson:"id"`
	Tag      int                `json:"tag,omitempty" bson:"tag"`
	Title    string             `json:"name" bson:"name"`
	Info     string             `json:"info" bson:"info"`
	Content  string             `json:"content,,omitempty" bson:"content"`
	SaveTime int64              `json:"save_time" bson:"save_time"`
}

type PromptConfig struct {
	Prompt   string
	Function string
}
