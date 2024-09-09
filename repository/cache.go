package repository

import (
	"ai-host/constant"
	"ai-host/global"
	"ai-host/model"
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

func (h *History) GetUserSession(ctx context.Context, msg model.RunRequest) (*model.RunRequest, error) {
	key := fmt.Sprintf(constant.UserCommit, strconv.Itoa(msg.Tag))
	date, err := h.Cache.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			msgJson, err := json.Marshal(msg)
			if err != nil {
				return nil, err
			}
			h.Cache.Set(ctx, key, string(msgJson), 5*time.Minute)
			return nil, nil
		}
	}
	var result *model.RunRequest
	err = json.Unmarshal([]byte(date), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
