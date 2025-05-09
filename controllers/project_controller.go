package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KHUTHON-404-not-found/khuthon-server/config"
	"github.com/KHUTHON-404-not-found/khuthon-server/models"
)

// CreateProject 새 프로젝트 생성
func CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 생성 시간 설정
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Init_todo(c, project.Plant)

	c.JSON(http.StatusCreated, project)
}

// GetProject 프로젝트 조회
func GetProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	if err := config.DB.Where("projectID = ?", id).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject 프로젝트 정보 업데이트
func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	// 프로젝트 존재 확인
	if err := config.DB.Where("projectID = ?", id).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// 요청 데이터 바인딩
	var input models.Project
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 업데이트 시간 설정
	input.UpdatedAt = time.Now()

	// 프로젝트 정보 업데이트
	if err := config.DB.Model(&project).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject 프로젝트 삭제
func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	// 프로젝트 존재 확인
	if err := config.DB.Where("projectID = ?", id).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// 프로젝트 삭제 (관련 일기와 사용자 관계 처리 필요)
	if err := config.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// GetProjectsByUser 프로젝트 아이디로 사용자 조회
func GetProjectsByUser(c *gin.Context) {
	userID := c.Param("user_id")
	var projects []models.Project

	if err := config.DB.Where("userID = ?", userID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(projects) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found for this project"})
		return
	}
	var projectIDs []int
	for _, project := range projects {
		projectIDs = append(projectIDs, project.ProjectID)
	}
	c.JSON(http.StatusOK, projectIDs)
}

// GetAllProjects 모든 프로젝트 조회
func GetAllProjects(c *gin.Context) {
	var projects []models.Project

	if err := config.DB.Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}
