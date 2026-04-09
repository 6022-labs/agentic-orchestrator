package models

import (
	"time"

	"gorm.io/gorm"
)

// Agent collection holds the schema definition for the AgentCollection entity.
type AgentCollection struct {
	Address   string     `gorm:"primaryKey;uniqueIndex;type:string;not null" json:"address"`
	Name      string     `gorm:"type:string;not null" json:"name"`
	Symbol    string     `gorm:"type:string;not null" json:"symbol"`
	UpdatedAt time.Time  `gorm:"type:timestamp;not null" json:"updatedAt"`
	CreatedAt time.Time  `gorm:"type:timestamp;not null" json:"createdAt"`
	DeletedAt *time.Time `gorm:"type:timestamp" json:"deletedAt"`
}

func (entity *AgentCollection) BeforeCreate(tx *gorm.DB) (err error) {
	entity.CreatedAt = time.Now()
	return
}

func (entity *AgentCollection) BeforeSave(tx *gorm.DB) (err error) {
	entity.UpdatedAt = time.Now()
	return
}
