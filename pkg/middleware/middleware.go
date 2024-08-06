package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	// instantiation
	logger := logrus.New()

	// Set output
	logger.Out = os.Stdout
	logger.SetFormatter(&logrus.TextFormatter{})

	// Set log level
	logger.SetLevel(logrus.DebugLevel)

	return func(c *gin.Context) {
		// Define the base route
		baseRoute := "/"
		if c.Request.RequestURI == baseRoute {
			c.Next()
			return
		}

		// Start time
		startTime := time.Now()

		// Process request
		c.Next()

		// End time
		endTime := time.Now()

		// Execution time
		latencyTime := endTime.Sub(startTime)

		// Request method
		reqMethod := c.Request.Method

		// Request routing
		reqURI := c.Request.RequestURI

		// Status code
		statusCode := c.Writer.Status()

		// Request IP
		clientIP := c.ClientIP()

		// User identifier
		userIdentifier := "-"

		// User ID
		userID := "-"

		// Current time in the specified format
		currentTime := time.Now().Format("02/Jan/2006:15:04:05 -0700")

		// Size of the object returned to the client
		responseSize := c.Writer.Size()

		// Log format
		logEntry := fmt.Sprintf("%s %s %s [%s] \"%s %s HTTP/1.0\" %d %d %v",
			clientIP, userIdentifier, userID, currentTime, reqMethod, reqURI, statusCode, responseSize, latencyTime)

		logger.Log(logrus.DebugLevel, logEntry)
	}
}
