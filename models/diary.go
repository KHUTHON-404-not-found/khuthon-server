package models

import (
	"time"
)

// Diary 모델 정의
type Diary struct {
	DiaryID   int       `json:"diary_id" gorm:"column:diaryID;primaryKey;autoIncrement"`
	ProjectID int       `json:"project_id" gorm:"column:projectID"`
	Title     string    `json:"title" gorm:"column:title"`
	Date      time.Time `json:"date" gorm:"column:date"`
	PhotoURL  string    `json:"photo_url" gorm:"column:photo_url"`
	Content   string    `json:"content" gorm:"column:content"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName 테이블 이름 지정
func (Diary) TableName() string {
	return "Diary"
}
