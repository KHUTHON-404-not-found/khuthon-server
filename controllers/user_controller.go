package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/KHUTHON-404-not-found/khuthon-server/config"
	"github.com/KHUTHON-404-not-found/khuthon-server/models"
	"github.com/KHUTHON-404-not-found/khuthon-server/utils"
)

// CreateUser 새 사용자 생성
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 비밀번호 해시화
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
	//	return
	//}
	//user.Password = string(hashedPassword)

	// 생성 시간 설정
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 비밀번호 필드 제외하고 응답
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// GetUser 사용자 조회
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.Where("userID = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 비밀번호 필드 제외하고 응답
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// UpdateUser 사용자 정보 업데이트
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// 사용자 존재 확인
	if err := config.DB.Where("userID = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 요청 데이터 바인딩
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 비밀번호가 변경된 경우 해시화
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		input.Password = string(hashedPassword)
	}

	// 업데이트 시간 설정
	input.UpdatedAt = time.Now()

	// 사용자 정보 업데이트
	if err := config.DB.Model(&user).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 비밀번호 필드 제외하고 응답
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// DeleteUser 사용자 삭제
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// 사용자 존재 확인
	if err := config.DB.Where("userID = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 사용자 삭제
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetUserByEmail 이메일로 사용자 조회
func GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	var user models.User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 비밀번호 필드 제외하고 응답
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// Login 사용자 로그인
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 비밀번호 검증
	if err := utils.VerifyPassword(user.Password, input.Password, c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// JWT 토큰 생성
	atPrivateKey := os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	rtPrivateKey := os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
	accessTTL := 60 * time.Minute
	refreshTTL := 7 * 24 * time.Hour
	accessToken, err := utils.CreateToken(accessTTL, user.UserID, atPrivateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token: " + err.Error()})
		return
	}
	refreshToken, err := utils.CreateToken(refreshTTL, user.UserID, rtPrivateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}
	// 쿠키 설정
	c.SetCookie("access_token", accessToken, int(accessTTL.Seconds()), "/", "", false, true)

	// 성공 응답
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken, "user": user})
}

// GetAllUsers 모든 사용자 조회 (관리자용)
func GetAllUsers(c *gin.Context) {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 비밀번호 필드 제외
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}
