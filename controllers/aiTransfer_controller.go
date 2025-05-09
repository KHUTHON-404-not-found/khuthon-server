package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 요청 바디 구조체
type MissionRequest struct {
	ImageURL string `json:"image_url" binding:"required"`
	Mission  string `json:"mission" binding:"required"`
}

// 프론트로부터 받은 URL과 mission을 외부로 전달
func TransferToMissionCheck(c *gin.Context) {
	var req MissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 형식"})
		return
	}

	payload := map[string]string{
		"image_url": req.ImageURL,
		"mission":   req.Mission,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON 변환 실패"})
		return
	}

	resp, err := http.Post("http://163.180.186.42:8000/checkMission", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "요청 전송 실패"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "응답 읽기 실패"})
		return
	}

	c.JSON(resp.StatusCode, gin.H{
		"result": string(body),
	})
}
