package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KHUTHON-404-not-found/khuthon-server/config"
	"github.com/KHUTHON-404-not-found/khuthon-server/models"
)

// CreateTodo 새 할일 생성
func CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 생성 시간 설정
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	if err := config.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// GetTodo 할일 조회
func GetTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	if err := config.DB.Where("todoID = ?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// UpdateTodo 할일 정보 업데이트
func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	// 할일 존재 확인
	if err := config.DB.Where("todoID = ?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	// 요청 데이터 바인딩
	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 업데이트 시간 설정
	input.UpdatedAt = time.Now()

	// 할일 정보 업데이트
	if err := config.DB.Model(&todo).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// DeleteTodo 할일 삭제
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	// 할일 존재 확인
	if err := config.DB.Where("todoID = ?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	// 할일 삭제
	if err := config.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}

// GetAllTodos 모든 할일 조회
func GetAllTodos(c *gin.Context) {
	var todos []models.Todo

	if err := config.DB.Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// GetTodosByDate 날짜별 할일 조회
func GetTodosByDate(c *gin.Context) {
	dateStr := c.Query("date")
	date, err := time.Parse("2006-01-02", dateStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	var todos []models.Todo

	// 날짜 범위로 검색 (하루)
	nextDay := date.AddDate(0, 0, 1)

	if err := config.DB.Where("date < ? AND complete = ?", nextDay, false).Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// CompleteTodo 할일 완료 상태 토글
func CompleteTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	// 할일 존재 확인
	if err := config.DB.Where("todoID = ?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	// 완료 상태 토글
	todo.Complete = !todo.Complete
	todo.UpdatedAt = time.Now()

	if err := config.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// GetTodosByProject Todo 아이디로 투두 조회
func GetTodosByProject(c *gin.Context) {
	projectID := c.Param("project_id")
	var todos []models.Todo

	if err := config.DB.Where("projectID = ?", projectID).Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}
