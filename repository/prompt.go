package repository

import (
	"ai-host/constant"
	"ai-host/global"
	"ai-host/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Prompt struct {
	promptCollection  *mongo.Collection
	historyCollection *mongo.Collection
}

type PromptInterface interface {
	// GetPromptList 查询prompt列表
	GetPromptList(ctx context.Context) (map[string]*model.Prompt, error)
	// GetFunctionList 查询所有function
	GetFunctionList(ctx context.Context) (map[string]*model.Prompt, error)
}

func NewPrompt() *Prompt {
	return &Prompt{
		promptCollection:  global.MongoClient.Collection(constant.AiHostPrompt),
		historyCollection: global.MongoClient.Collection(constant.AiTestHistory),
	}
}

func (p *Prompt) GetPromptList(ctx context.Context) (map[string]*model.Prompt, error) {
	var list []*model.Prompt
	cursor, err := p.promptCollection.Find(ctx, bson.D{{"version", bson.M{"$ne": constant.AiFunctions}}})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	result := make(map[string]*model.Prompt)
	for _, v := range list {
		result[v.Version] = v
	}
	return result, nil
}

func (p *Prompt) GetFunction(ctx context.Context) (*model.Prompt, error) {
	var functions model.Prompt
	err := p.promptCollection.FindOne(ctx, bson.D{{"version", constant.AiFunctions}}).Decode(&functions)
	if err != nil {
		return nil, err
	}
	return &functions, nil
}

func (p *Prompt) GetHistoryList(ctx context.Context) ([]*model.History, error) {
	var list []*model.History
	projection := bson.D{{"_id", 1}, {"name", 1}}
	// 创建查询选项
	findOptions := options.Find().SetProjection(projection)
	cursor, err := p.promptCollection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (p *Prompt) InsertHistory(ctx context.Context, req model.SaveRequest, content string) error {
	_, err := p.historyCollection.InsertOne(ctx, model.History{
		Tag:      req.Tag,
		Title:    req.Title,
		Info:     req.Info,
		Content:  content,
		SaveTime: time.Now().UnixMilli(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Prompt) GetHistoryById(ctx context.Context, id primitive.ObjectID) (*model.History, error) {
	var history *model.History
	err := p.historyCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(history)
	if err != nil {
		return nil, err
	}
	return history, nil
}

// 插入Prompt
func (p *Prompt) InsertPrompt(ctx context.Context, data *model.Prompt) error {
	_, err := p.historyCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// 删除历史记录
func (p *Prompt) DeleteHistory(ctx context.Context, id primitive.ObjectID) error {
	_, err := p.historyCollection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}
