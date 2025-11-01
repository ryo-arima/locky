package middleware

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/config"
)

func LoggerWithConfig(conf config.BaseConfig) gin.HandlerFunc {
	lv := strings.ToLower(strings.TrimSpace(conf.YamlConfig.Application.Server.LogLevel))
	if lv == "" {
		lv = "info"
	}
	logger := log.New(os.Stdout, "", log.LstdFlags)
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		msg := "Request processed"
		if len(c.Errors) > 0 {
			msg = c.Errors.Last().Error()
		}

		// level mapping
		level := "INFO"
		if status >= 500 {
			level = "ERROR"
		} else if status >= 400 {
			level = "WARN"
		}

		// Output determination
		allow := func(outLevel string) bool {
			order := map[string]int{"debug": 0, "info": 1, "warn": 2, "error": 3}
			return order[strings.ToLower(outLevel)] >= order[lv]
		}
		if !allow(level) && !(lv == "debug" && level == "INFO") {
			return
		}

		if lv == "debug" {
			logger.Printf("[%s] %d %s %s %s %v UA=%s HDR=%v ERR=%v", level, status, method, path, clientIP, latency, c.Request.UserAgent(), c.Request.Header, c.Errors)
		} else {
			logger.Printf("[%s] %d %s %s %s %v %s", level, status, method, path, clientIP, latency, msg)
		}
	}
}
