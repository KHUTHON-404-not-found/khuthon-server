package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KHUTHON-404-not-found/khuthon-server/config"
	"github.com/KHUTHON-404-not-found/khuthon-server/models"
)

// CreateDiary 새 일기 생성
func CreateDiary(c *gin.Context) {
	var diary models.Diary
	if err := c.ShouldBindJSON(&diary); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 생성 시간 설정
	diary.CreatedAt = time.Now()
	diary.UpdatedAt = time.Now()

	if err := config.DB.Create(&diary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, diary)
}

// GetDiary 일기 조회
func GetDiary(c *gin.Context) {
	id := c.Param("id")
	var diary models.Diary

	if err := config.DB.Where("diaryID = ?", id).First(&diary).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
		return
	}

	c.JSON(http.StatusOK, diary)
}

// UpdateDiary 일기 정보 업데이트
func UpdateDiary(c *gin.Context) {
	id := c.Param("id")
	var diary models.Diary

	// 일기 존재 확인
	if err := config.DB.Where("diaryID = ?", id).First(&diary).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
		return
	}

	// 요청 데이터 바인딩
	var input models.Diary
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 업데이트 시간 설정
	input.UpdatedAt = time.Now()

	// 일기 정보 업데이트
	if err := config.DB.Model(&diary).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, diary)
}

// DeleteDiary 일기 삭제
func DeleteDiary(c *gin.Context) {
	id := c.Param("id")
	var diary models.Diary

	// 일기 존재 확인
	if err := config.DB.Where("diaryID = ?", id).First(&diary).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
		return
	}

	// 일기 삭제
	if err := config.DB.Delete(&diary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diary deleted successfully"})
}

// GetAllDiaries 모든 일기 조회
func GetAllDiaries(c *gin.Context) {
	var diaries []models.Diary

	if err := config.DB.Find(&diaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, diaries)
}

// GetDiariesByProject 프로젝트별 일기 조회
func GetDiariesByProject(c *gin.Context) {
	projectID := c.Param("project_id")
	var diaries []models.Diary

	if err := config.DB.Where("projectID = ?", projectID).Find(&diaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, diaries)
}

// GetDiariesByDate 날짜별 일기 조회
func GetDiariesByDate(c *gin.Context) {
	dateStr := c.Query("date")
	date, err := time.Parse("2006-01-02", dateStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	var diaries []models.Diary

	// 날짜 범위로 검색 (하루)
	nextDay := date.AddDate(0, 0, 1)

	if err := config.DB.Where("date >= ? AND date < ?", date, nextDay).Find(&diaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, diaries)
}

// UploadDiaryPhoto 일기 사진 업로드 (실제 구현시 파일 업로드 처리 필요)
func UpdateDiaryPhotoURL(c *gin.Context) {
	id := c.Param("id")
	var diary models.Diary

	// 기존 일기 조회
	if err := config.DB.Where("diaryID = ?", id).First(&diary).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
		return
	}

	// photo_url만 받기
	var input struct {
		PhotoURL string `json:"photo_url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 업데이트
	diary.PhotoURL = input.PhotoURL
	diary.UpdatedAt = time.Now()

	if err := config.DB.Save(&diary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, diary)
}