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

	weeksRaw, ok := responseBody["weeks"]
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
		weekNumber := weekData["week"].(float64)
		todoList := weekData["tasks"].([]interface{})
		wateringFrequency := weekData["watering"].(map[string]interface{})["frequency"].(float64)

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
			log.Println("todoReq", todoReq)
			createTodoFromModel(c, todoReq)
		}

		for i := 0; i < 7/int(wateringFrequency); i++ {
			wateringReq := models.Todo{
				ProjectID: projectID,
				Content:   "물 주기",
				Date:      waterday,
				Complete:  false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			waterday = waterday.AddDate(0, 0, int(wateringFrequency))
			createTodoFromModel(c, wateringReq)
		}
	}

	c.JSON(http.StatusOK, gin.H{"body": responseBody})

}

func createTodoFromModel(c *gin.Context, todoReq models.Todo) {
	// 생성 시간 설정
	todoReq.CreatedAt = time.Now()
	todoReq.UpdatedAt = time.Now()

	if err := config.DB.Create(&todoReq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todoReq)
}
