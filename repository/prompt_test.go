package repository

import (
	"ai-host/global"
	"ai-host/model"
	"context"
	"fmt"
	"testing"
)

func newConfig() {
	global.Config = model.Config{
		Mongo: model.MongoConf{
			Address:     "mongodb://localhost:27017/wooplus?retryWrites=true",
			MaxPoolSize: 10,
		},
	}
	// 初始化依赖
	global.InitMongo()
}

func TestInsertPrompt(t *testing.T) {
	data := &model.Prompt{
		Version:            "functions",
		System:             `{"name":"next_step","description":"What should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, end,or banter.e.g.next_step:introduction","strict":false,"parameters":{"type":"object","properties":{"next_step":{"description":"What should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, or banter","type":"string"}},"required":["next_step"]}}`,
		Introduction:       `{"name":"introduction","description":"The content of introduction and what should we do next step, including introduction, icebreaker question, user_answer, follow-up question, new question, reply,end, or banter","strict":false,"parameters":{"type":"object","properties":{"introduction":{"description":"The content of brief interesting and attractive introduction which should be based on both profiles.Keep introductions brief, aiming to complete them in two to three sentences.","type":"string"},"next_step":{"description":"What should we do next step, including introduction, icebreaker question, user_answer, follow-up question, new question, reply, or banter","type":"string"}},"required":["introduction","next_step"]}}`,
		IcebreakerQuestion: `{"name":"icebreaker_question","description":"The content of icebreaker_question and what should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, or banter","strict":false,"parameters":{"type":"object","properties":{"icebreaker_question":{"description":"The content of icebreaker_question which should be based on both profiles. Create a multiple-choice question with only two options, where the question and content are related to dating and are as relevant as possible to both parties.Two options with content only, without the use of labels like A or B.","type":"string"},"A":{"description":"The first option for the multiple-choice question.","type":"string"},"B":{"description":"The second option for the multiple-choice question.","type":"string"},"next_step":{"description":"What should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, or banter.","type":"string"}},"required":["icebreaker_question","A","B","next_step"]}}`,
		FollowUpQuestion:   `{"name":"follow_up_question","description":"The content of follow_up_question and what should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, end,or banter","strict":false,"parameters":{"type":"object","properties":{"follow_up_question":{"description":"The content of follow_up_question which is an additional question asked to gain further clarification or explore a topic in more depth after the initial question. Do not ask separate questions for each person","type":"string"},"next_step":{"description":"What should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, or banter","type":"string"}},"required":["follow_up_question","next_step"]}}`,
		NewQuestion:        `{"name":"new_question","description":"The content of new_question and what should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, end,or banter","strict":false,"parameters":{"type":"object","properties":{"new_question":{"description":"The content of the new_questions is related to the dating theme. The purpose is to allow both parties to increase understanding and raise interesting and novel issues that can be discussed at a high level.","type":"string"},"next_step":{"description":"What should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, or banter","type":"string"}},"required":["new_question","next_step"]}}`,
		Banter:             `{"name":"banter","description":"The content of banter and what should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, end,or banter","strict":false,"parameters":{"type":"object","properties":{"banter":{"description":"The content of banter which is to analyze the content of both parties' responses to find common ground and use that for banter. The banter should be very short in one sentence","type":"string"},"next_step":{"description":"What should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, or banter","type":"string"}},"required":["banter","next_step"]}}`,
		Reply:              `{"name":"reply","description":"The content of reply and what should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, end,or banter","strict":false,"parameters":{"type":"object","properties":{"reply":{"description":"The content of reply which is a response to users to make the dating atmosphere more relaxed,Don’t ask questions in reply","type":"string"},"next_step":{"description":"What should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, or banter","type":"string"}},"required":["reply","next_step"]}}`,
		WarnUp:             `{"name":"warm_up","description":"The content of warm_up and what should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, end,or banter","strict":false,"parameters":{"type":"object","properties":{"warm_up":{"description":"The content of warm_up which remind the two parties who have not communicated with each other for some time that a moderator is needed to break the deadlock. You can use a humorous tone to remind both parties to continue chatting, or you can ask follow_up_questions to extend the previous topic.  ","type":"string"},"next_step":{"description":"What should we do next step, including introduction, icebreaker_question, user_answer, follow_up_question, new_question, reply, or banter","type":"string"}},"required":["warm_up","next_step"]}}`,
		UserAnswer:         `{"name":"user_answer","description":"what should we do next step after answer,including follow_up_question, new_question, reply, final_round,or banter","strict":false,"parameters":{"type":"object","properties":{"next_step":{"description":"Please based on their answers, tell me What should we do next step after answer. analyze their responses， If the answer has a follow-up value, you give a type of follow_up_question based on the user_answer. If the user_answer has no follow-up value, you give a type of banter before giving a type of new_question. If the user's response contains disappointment, dissatisfaction, or similar negative emotions to the other one, or if the user directly addresses you, you need to give a type of reply.If the user expresses their intention to end the game, you need to provide the type of final_round","type":"string"}},"required":["next_step"]}}`,
		FinalRound:         `{"name":"final_round","description":"The content of final_round and what should we do next step, including end","strict":false,"parameters":{"type":"object","properties":{"final_round":{"description":"The content of final_round is, when the user's willingness to end is not strong, try to restore the user. When the user has a strong desire to end, appease the user's emotions and give closing words. The content in this section needs to be lighthearted and humorous ","type":"string"},"next_step":{"description":"What should we do next step, including end","type":"string"}},"required":[" final_round","next_step"]}}`,
	}

	newConfig()

	ctx := context.Background()
	err := NewPrompt().InsertPrompt(ctx, data)
	if err != nil {
		fmt.Println(err)
	}
}
