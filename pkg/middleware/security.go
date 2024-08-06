package middleware

import (
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Throttle() gin.HandlerFunc {
	var (
		requestPerSecond float64 = 7
		serverConfig             = config.GetConfig().Server
	)

	if serverConfig.RequestPerSecond != 0 {
		requestPerSecond = serverConfig.RequestPerSecond
	}

	lmt := tollbooth.NewLimiter(requestPerSecond, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"X-Forwarded-For", "X-Real-IP", "RemoteAddr"})
	lmt.SetMethods([]string{"GET", "POST", "PUT", "DELETE"})

	return func(c *gin.Context) {
		if isExemptIP(c.ClientIP(), serverConfig.ExemptFromThrottle) {
			c.Next()
		} else {
			httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
			if httpError != nil {
				c.AbortWithStatusJSON(httpError.StatusCode, utility.BuildErrorResponse(httpError.StatusCode, "error", lmt.GetMessage(), httpError.Message, nil))
			} else {
				c.Next()
			}
		}
	}
}

func isExemptIP(ip string, exemptIPs []string) bool {
	for _, exemptIP := range exemptIPs {
		if ip == exemptIP {
			return true
		}
	}
	return false
}

// Security middleware
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-XSS-Protection
		c.Writer.Header().Add("X-XSS-Protection", "1; mode=block")

		// HTTP Strict Transport Security
		c.Writer.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// X-Frame-Options
		c.Writer.Header().Add("X-Frame-Options", "SAMEORIGIN")

		// X-Content-Type-Options
		c.Writer.Header().Add("X-Content-Type-Options", "nosniff")

		// Content Security Policy
		c.Writer.Header().Add("Content-Security-Policy", "default-src 'self';")

		// X-Permitted-Cross-Domain-Policies
		c.Writer.Header().Add("X-Permitted-Cross-Domain-Policies", "none")

		// Referrer-Policy
		c.Writer.Header().Add("Referrer-Policy", "no-referrer")

		// Feature-Policy
		c.Writer.Header().Add("Feature-Policy", "microphone 'none'; camera 'none'")

		c.Next()
	}
}
