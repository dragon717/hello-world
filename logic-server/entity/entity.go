package entity

import (
	"encoding/json"
)

// EntityType 定义实体类型
type EntityType string

const (
	EntityTypeSystem EntityType = "system"
	EntityTypeUser   EntityType = "user"
	EntityTypeNPC    EntityType = "npc"
)

// Entity 基础实体接口
type Entity interface {
	GetID() string
	GetType() EntityType
	GetData() map[string]interface{}
	SetData(key string, value interface{})
	ToJSON() ([]byte, error)
}

// BaseEntity 基础实体结构
type BaseEntity struct {
	ID   string                 `json:"id"`
	Type EntityType             `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// NewBaseEntity 创建基础实体
func NewBaseEntity(id string, entityType EntityType) *BaseEntity {
	return &BaseEntity{
		ID:   id,
		Type: entityType,
		Data: make(map[string]interface{}),
	}
}

func (e *BaseEntity) GetID() string {
	return e.ID
}

func (e *BaseEntity) GetType() EntityType {
	return e.Type
}

func (e *BaseEntity) GetData() map[string]interface{} {
	return e.Data
}

func (e *BaseEntity) SetData(key string, value interface{}) {
	e.Data[key] = value
}

func (e *BaseEntity) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// SystemEntity 系统级实体
type SystemEntity struct {
	*BaseEntity
}

func NewSystemEntity(id string) *SystemEntity {
	return &SystemEntity{
		BaseEntity: NewBaseEntity(id, EntityTypeSystem),
	}
}

// UserEntity 用户级实体
type UserEntity struct {
	*BaseEntity
	Username string `json:"username"`
}

func NewUserEntity(id string, username string) *UserEntity {
	return &UserEntity{
		BaseEntity: NewBaseEntity(id, EntityTypeUser),
		Username:   username,
	}
}

// NPCEntity NPC实体
type NPCEntity struct {
	*BaseEntity
	Name string `json:"name"`
}

func NewNPCEntity(id string, name string) *NPCEntity {
	return &NPCEntity{
		BaseEntity: NewBaseEntity(id, EntityTypeNPC),
		Name:       name,
	}
}
