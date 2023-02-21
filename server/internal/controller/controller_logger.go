package controller

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func loggerMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()

		context.Next()

		lat := time.Now().Sub(start)

		path := context.Request.URL.Path
		method := context.Request.Method
		status := context.Writer.Status()
		bodySize := context.Writer.Size()
		clientIp := context.ClientIP()

		l := log.WithFields(log.Fields{
			"status":   status,
			"latency":  lat,
			"client":   clientIp,
			"bodySize": bodySize,
		})

		errs := context.Errors.ByType(gin.ErrorTypePrivate)

		httpUrlStr := fmt.Sprintf("%s %s", method, path)

		if len(errs) != 0 {
			httpUrl := color.New(color.FgHiRed).Sprint(httpUrlStr)
			l.Errorf("handled %s with errors: %s", httpUrl, errs)
			return
		}

		if status >= 400 && status < 499 {
			httpUrl := color.New(color.FgHiYellow).Sprint(httpUrlStr)
			l.Warnf("handled %s with client error", httpUrl)
			return
		}

		if status >= 500 {
			httpUrl := color.New(color.FgHiRed).Sprint(httpUrlStr)
			l.Errorf("handled %s with server error", httpUrl)
			return
		}

		httpUrl := color.New(color.FgHiGreen).Sprint(httpUrlStr)
		l.Infof("handled %s successfully", httpUrl)
	}
}
