package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hello-world/logic-server/entity"
	"github.com/redis/go-redis/v9"
)

// Storage 存储接口
type Storage interface {
	SaveEntity(ctx context.Context, e entity.Entity) error
	GetEntity(ctx context.Context, id string) (entity.Entity, error)
	DeleteEntity(ctx context.Context, id string) error
}

// RedisStorage Redis存储实现
type RedisStorage struct {
	client *redis.Client
}

// NewRedisStorage 创建Redis存储实例
func NewRedisStorage(host string, port int, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStorage{client: client}, nil
}

// SaveEntity 保存实体
func (s *RedisStorage) SaveEntity(ctx context.Context, e entity.Entity) error {
	data, err := e.ToJSON()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("entity:%s", e.GetID())
	return s.client.Set(ctx, key, data, 0).Err()
}

// GetEntity 获取实体
func (s *RedisStorage) GetEntity(ctx context.Context, id string) (entity.Entity, error) {
	key := fmt.Sprintf("entity:%s", id)
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	// 解析基础实体以确定类型
	var base entity.BaseEntity
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	// 根据类型创建对应的实体
	switch base.Type {
	case entity.EntityTypeSystem:
		var system entity.SystemEntity
		if err := json.Unmarshal(data, &system); err != nil {
			return nil, err
		}
		return &system, nil
	case entity.EntityTypeUser:
		var user entity.UserEntity
		if err := json.Unmarshal(data, &user); err != nil {
			return nil, err
		}
		return &user, nil
	case entity.EntityTypeNPC:
		var npc entity.NPCEntity
		if err := json.Unmarshal(data, &npc); err != nil {
			return nil, err
		}
		return &npc, nil
	default:
		return nil, fmt.Errorf("unknown entity type: %s", base.Type)
	}
}

// DeleteEntity 删除实体
func (s *RedisStorage) DeleteEntity(ctx context.Context, id string) error {
	key := fmt.Sprintf("entity:%s", id)
	return s.client.Del(ctx, key).Err()
}
