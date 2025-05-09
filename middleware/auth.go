package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/KHUTHON-404-not-found/khuthon-server/config"
	"github.com/KHUTHON-404-not-found/khuthon-server/models"
)

// JWTAuth JWT 인증 미들웨어
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &jwt.StandardClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 실제 구현에서는 환경 변수나 설정 파일에서 시크릿 키를 가져와야 함
			return []byte("your_jwt_secret_key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 사용자 ID를 컨텍스트에 저장
		c.Set("userID", claims.Subject)
		c.Next()
	}
}

// GenerateToken JWT 토큰 생성
func GenerateToken(userID string) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), // 토큰 유효 기간 3일
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 실제 구현에서는 환경 변수나 설정 파일에서 시크릿 키를 가져와야 함
	return token.SignedString([]byte("your_jwt_secret_key"))
}

// LoginHandler 로그인 후 JWT 토큰 발급
func LoginHandler(c *gin.Context) {
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

	// 비밀번호 검증 (실제 구현시 bcrypt 등으로 해시 비교 필요)
	// bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	// JWT 토큰 생성
	token, err := GenerateToken(string(user.UserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.UserID,
			"email": user.Email,
		},
	})
}
