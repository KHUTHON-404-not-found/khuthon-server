package utils

import (
	"crypto/subtle"
	"fmt"

	"github.com/gin-gonic/gin"
)

func VerifyPassword(hashedPassword string, candidatePassword string, c *gin.Context) error {
	if subtle.ConstantTimeCompare([]byte(candidatePassword), []byte(hashedPassword)) != 1 {
		err := fmt.Errorf("Invalid email or password")
		return err
	}
	return nil
}
