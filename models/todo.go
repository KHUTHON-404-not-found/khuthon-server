package models

import (
	"time"
)

// Todo 모델 정의
type Todo struct {
	TodoID    int       `json:"todo_id" gorm:"column:todoID;primaryKey;autoIncrement"`
	Date      time.Time `json:"date" gorm:"column:date"`
	Content   string    `json:"content" gorm:"column:content"`
	Complete  bool      `json:"complete" gorm:"column:complete"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName 테이블 이름 지정
func (Todo) TableName() string {
	return "Todo"
}
