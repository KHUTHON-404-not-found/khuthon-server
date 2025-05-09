package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init_todo(c *gin.Context, plant string) {
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

	// 응답 리턴 (필요 시 처리 가능)
	c.JSON(http.StatusOK, gin.H{"message": "Request sent successfully", "Body": resp.Body})
}
