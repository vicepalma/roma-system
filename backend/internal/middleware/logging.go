package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderRequestID = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(HeaderRequestID)
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Writer.Header().Set(HeaderRequestID, rid)
		c.Set(HeaderRequestID, rid)
		c.Next()
	}
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf(
			`level=info ts=%s rid=%s method=%s path=%s status=%d dur_ms=%d ua="%s"`,
			time.Now().Format(time.RFC3339),
			c.Writer.Header().Get(HeaderRequestID),
			c.Request.Method,
			c.FullPath(),
			c.Writer.Status(),
			time.Since(start).Milliseconds(),
			c.Request.UserAgent(),
		)
	}
}
