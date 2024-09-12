package repository

import (
	"ai-dating/constant"
	"ai-dating/global"
	"ai-dating/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type History struct {
	Cache *redis.Client
}

func NewCacheHistory() *History {
	return &History{
		Cache: global.RedisClint,
	}
}

type IHistory interface {
	// GetSessionHistory 查询会话历史
	GetSessionHistory(ctx context.Context, tag string) ([]model.OpenAIMessageType, error)
	// AddSessionHistory 加入会话历史
	AddSessionHistory(ctx context.Context, tag string, msg []model.OpenAIMessageType) error
	// AddUserProfile  用户profile加入缓存
	AddUserProfile(ctx context.Context, msg model.StartRequest) error
	// GetUserProfile 查询用户的profile
	GetUserProfile(ctx context.Context, tag int) (*model.Profile, error)
}

func (h *History) GetSessionHistory(ctx context.Context, tag string) ([]model.OpenAIMessageType, error) {
	var result []model.OpenAIMessageType
	key := fmt.Sprintf(constant.AIHostSessionRecord, tag)
	date, err := h.Cache.Get(ctx, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, err
		} else {
			return result, nil
		}
	}

	err = json.Unmarshal([]byte(date), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (h *History) AddSessionHistory(ctx context.Context, tag string, msg []model.OpenAIMessageType) error {
	key := fmt.Sprintf(constant.AIHostSessionRecord, tag)
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = h.Cache.Set(ctx, key, string(msgJson), 30*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

// Profile缓存
func (h *History) GetUserProfile(ctx context.Context, tag int) (*model.Profile, error) {
	key := fmt.Sprintf(constant.UserCommit, strconv.Itoa(tag))
	date, err := h.Cache.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New(" GetUserProfile is nil! ")
		}
	}
	var result *model.Profile
	err = json.Unmarshal([]byte(date), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AddProfileCache
func (h *History) AddUserProfile(ctx context.Context, msg model.StartRequest) error {
	key := fmt.Sprintf(constant.UserCommit, strconv.Itoa(msg.Tag))
	msgJson, err := json.Marshal(model.Profile{
		Male:   msg.Male,
		Female: msg.Female,
	})
	if err != nil {
		return err
	}
	err = h.Cache.Set(ctx, key, string(msgJson), 30*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
