package monitoring

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"time"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("API request | method: %s | path: %s | status_code: %d | duration: %s", method, path, statusCode, duration)
	}
}

func RequestBodyLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		reader := ioutil.NopCloser(bytes.NewReader(buf))
		c.Request.Body = reader
		logrus.Info(string(buf))
		c.Next()
	}
}
