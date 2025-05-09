package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/KHUTHON-404-not-found/khuthon-server/config"
	"github.com/KHUTHON-404-not-found/khuthon-server/models"
	"github.com/gin-gonic/gin"
)

func Init_todo(c *gin.Context, plant string, projectID int) {
	var req struct {
		Keyword string `json:"keyword"`
	}
	req.Keyword = plant

	payload, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return
	}

	// HTTP POST 요청
	resp, err := http.Post("http://163.180.186.42:8000/initCrop", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to crop server"})
		return
	}
	defer resp.Body.Close()

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"body": responseBody})
	bodyRaw, ok := responseBody["body"]
	if !ok {
		log.Println("body key not found in responseBody")
		return
	}

	body, ok := bodyRaw.(map[string]interface{})
	if !ok {
		log.Println("body is not a map[string]interface{}")
		return
	}

	weeksRaw, ok := body["weeks"]
	if !ok {
		log.Println("weeks key not found in body")
		return
	}

	weeksList, ok := weeksRaw.([]interface{})
	if !ok {
		log.Println("weeks is not a list")
		return
	}

	today := time.Now()
	waterday := today
	for _, week := range weeksList {
		weekData := week.(map[string]interface{})
		weekNumber := weekData["week"].(int)
		todoList := weekData["tasks"].([]interface{})
		wateringFrequency := weekData["watering"].(map[string]interface{})["frequency"].(int)

		weekStartDate := today.AddDate(0, 0, int(weekNumber-1)*7)

		for _, todo := range todoList {
			//createTodo(c)
			todoReq := models.Todo{
				ProjectID: projectID,
				Content:   todo.(string),
				Date:      weekStartDate,
				Complete:  false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			createTodoFromModel(c, todoReq)
		}

		for i := 0; i < 7/wateringFrequency; i++ {
			wateringReq := models.Todo{
				ProjectID: projectID,
				Content:   "물 주기",
				Date:      waterday,
				Complete:  false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			waterday = today.AddDate(0, 0, wateringFrequency)
			createTodoFromModel(c, wateringReq)
		}
	}

	c.JSON(http.StatusOK, gin.H{"body": responseBody})

}

func createTodoFromModel(c *gin.Context, todoReq models.Todo) {
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
