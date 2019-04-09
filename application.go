/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush LoggerWriter plugin]
 */

package logger

import (
	"fmt"
	"path"
	"time"

	"github.com/2637309949/bulrush"
	"github.com/gin-gonic/gin"
)

// LOGLEVEL -
type LOGLEVEL int

const (
	// SYSLEVEL  -
	SYSLEVEL LOGLEVEL = 0 + iota
	// USERLEVEL -
	USERLEVEL
	// SYSSTROBE -
	SYSSTROBE = 1024 * 1024 * 5
	// USERSTROBE -
	USERSTROBE = 1024 * 1024 * 5
)

type (
	// LoggerWriter plugin
	LoggerWriter struct {
		bulrush.PNBase
		cfg *bulrush.Config
	}
)

// Plugin for Recovery
func (logger *LoggerWriter) Plugin() bulrush.PNRet {
	return func(cfg *bulrush.Config, router *gin.RouterGroup) {
		logger.cfg = cfg
		router.Use(func(c *gin.Context) {
			fmt.Println("####  logger")
			logsDir := Some(LeftV(cfg.String("logs")), "logs").(string)
			logsDir = path.Join(".", logsDir)
			logPath := getLogFile(SYSLEVEL, logsDir)
			writer := createLog(logPath)
			start := time.Now()
			path := c.Request.URL.Path
			raw := c.Request.URL.RawQuery
			if raw != "" {
				path = path + "?" + raw
			}
			clientIP := c.ClientIP()
			method := c.Request.Method

			fmt.Fprintf(writer, "[%v] => %s %6s %s\n", start.Format("2006/01/02 15:04:05"), clientIP, method, path)
			c.Next()
			end := time.Now()
			latency := float64(end.Sub(start) / time.Millisecond)
			fmt.Fprintf(writer, "[%v] <= %.2fms %s %6s %s\n", end.Format("2006/01/02 15:04:05"), latency, clientIP, method, path)
		})
	}
}

// Writer -
// fileName start with "1"
// User level
func (logger *LoggerWriter) Writer(info string) {
	logsDir := Some(LeftV(logger.cfg.String("logs")), "logs").(string)
	logsDir = path.Join(".", logsDir)
	logPath := getLogFile(USERLEVEL, logsDir)
	writer := createLog(logPath)
	out := writer
	start := time.Now()
	fmt.Fprintf(out, "[%v]<-S-> %s\n", start.Format("2006/01/02 15:04:05"), info)
}
