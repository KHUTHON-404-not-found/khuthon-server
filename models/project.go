package models

import (
	"time"
)

// Project 모델 정의
type Project struct {
	ProjectID int       `json:"project_id" gorm:"column:projectID;primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" gorm:"column:userID"`
	Plant     string    `json:"plant" gorm:"column:plant"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName 테이블 이름 지정
func (Project) TableName() string {
	return "Project"
}
