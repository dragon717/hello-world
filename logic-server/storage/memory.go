package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hello-world/logic-server/entity"
)

// MemoryStorage 内存存储实现
type MemoryStorage struct {
	mu    sync.RWMutex
	store map[string][]byte
}

// NewMemoryStorage 创建内存存储实例
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		store: make(map[string][]byte),
	}
}

// SaveEntity 保存实体
func (s *MemoryStorage) SaveEntity(ctx context.Context, e entity.Entity) error {
	data, err := e.ToJSON()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("entity:%s", e.GetID())
	s.store[key] = data
	return nil
}

// GetEntity 获取实体
func (s *MemoryStorage) GetEntity(ctx context.Context, id string) (entity.Entity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := fmt.Sprintf("entity:%s", id)
	data, ok := s.store[key]
	if !ok {
		return nil, fmt.Errorf("entity not found: %s", id)
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
func (s *MemoryStorage) DeleteEntity(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("entity:%s", id)
	delete(s.store, key)
	return nil
}
