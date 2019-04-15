/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush LoggerWriter plugin]
 */

package logger

import (
	"path"
	"time"

	"github.com/2637309949/bulrush"
	journal "github.com/2637309949/bulrush-addition/logger"
	"github.com/gin-gonic/gin"
)

type (
	// Logger plugin
	Logger struct {
		bulrush.PNBase
		cfg *bulrush.Config
	}
)

// Plugin for Recovery
func (logger *Logger) Plugin() bulrush.PNRet {
	return func(cfg *bulrush.Config, router *gin.RouterGroup) {
		logger.cfg = cfg
		journal := journal.CreateHTTPLogger(path.Join(".", Some(LeftV(cfg.String("logs")), "logs").(string)))
		router.Use(func(c *gin.Context) {
			start := time.Now()
			path := c.Request.URL.Path
			raw := c.Request.URL.RawQuery
			if raw != "" {
				path = path + "?" + raw
			}
			clientIP := c.ClientIP()
			method := c.Request.Method
			journal.Info("[%v] => %s %6s %s\n", start.Format("2006/01/02 15:04:05"), clientIP, method, path)
			c.Next()
			end := time.Now()
			latency := float64(end.Sub(start) / time.Millisecond)
			journal.Info("[%v] <= %.2fms %s %6s %s\n", end.Format("2006/01/02 15:04:05"), latency, clientIP, method, path)
		})
	}
}
