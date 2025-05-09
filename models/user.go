package models

import (
	"time"
)

// User 모델 정의
type User struct {
	UserID    int       `json:"user_id" gorm:"column:userID;primaryKey;autoIncrement"`
	ProjectID int       `json:"project_id" gorm:"column:projectID"`
	Email     string    `json:"email" gorm:"column:email;unique"`
	Password  string    `json:"password" gorm:"column:password"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName 테이블 이름 지정
func (User) TableName() string {
	return "User"
}
