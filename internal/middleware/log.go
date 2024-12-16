package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"kn-assignment/internal/log"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		method := c.Request.Method
		url := c.Request.URL

		logFormatRequest := fmt.Sprintf("Request API [%s]: -> Method: %s Total Time: %s", url, method, latency)
		LogJsonWithCtx(ctx, logFormatRequest, nil)
	}
}

func ResponseLogger(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Next()

		method := c.Request.Method
		url := c.Request.URL
		status := c.Writer.Status()

		logFormatResponse := fmt.Sprintf("Response API [%s]: -> Method: %s Status: %s", url, method, fmt.Sprint(status))
		LogJsonWithCtx(ctx, logFormatResponse, nil)
	}
}

func LogJsonWithCtx(ctx context.Context, format string, message any) {

	if message == nil {
		log.Info(ctx, format)
		return
	}

	if reflect.TypeOf(message).Kind() == reflect.String {
		log.Infof(ctx, format, message)
		return
	}

	if loggableByte, err := json.Marshal(message); err == nil {
		log.Infof(ctx, format, loggableByte)
	}
}
