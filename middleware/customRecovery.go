package middleware

import (
	"errors"
	"github.com/coderconquerer/go-login-app/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/utils"
	"log"
	"os"
	"runtime/debug"
)

func CustomRecovery() gin.HandlerFunc {
	// Open or create log file
	file, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}
	logger := log.New(file, "", log.LstdFlags)

	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Log the error and stack trace
				logger.Printf("[PANIC RECOVER] %v\n", r)
				logger.Printf("[STACK TRACE]\n%s", debug.Stack())

				// Optional: You can also return a JSON response
				c.AbortWithStatusJSON(500, common.NewInternalSeverErrorResponse(errors.New("module Server Error"), "Something went wrong with server", utils.ToString(debug.Stack())))
			}
		}()

		c.Next()
	}
}
