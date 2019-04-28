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
	journal "github.com/2637309949/bulrush-addition/logger"
	"github.com/gin-gonic/gin"
)

type (
	// Type logger type
	Type int
	// Payload in request
	Payload struct {
		StartUnix int64
		EndUnix   int64
		Type      Type
		Latency   float64
		Method    string
		IP        string
		URL       string
	}

	// Logger plugin
	Logger struct {
		bulrush.PNBase
		Format func(*Payload, *gin.Context) string
		Path   string
	}
)

const (
	// INT logger type
	INT Type = iota + 1
	// OUT logger type
	OUT
)

func defaultFormat(p *Payload, ctx *gin.Context) string {
	if p.Type == INT {
		startTime := time.Unix(p.StartUnix, 0).Format("2006/01/02 15:04:05")
		return fmt.Sprintf("[%v] => %s %6s %s", startTime, p.IP, p.Method, p.URL)
	} else if p.Type == OUT {
		endOfTime := time.Unix(p.EndUnix, 0).Format("2006/01/02 15:04:05")
		latency := float64(time.Unix(p.EndUnix, 0).Sub(time.Unix(p.StartUnix, 0)) / time.Millisecond)
		return fmt.Sprintf("[%v] <= %.2fms %s %6s %s", endOfTime, latency, p.IP, p.Method, p.URL)
	}
	return "NO FORMAT"
}

// Plugin for Recovery
func (logger *Logger) Plugin() bulrush.PNRet {
	return func(cfg *bulrush.Config, router *gin.RouterGroup) {
		logger.Path = Some(logger.Path, LeftV(cfg.String("logs"))).(string)
		journal := journal.CreateHTTPLogger(path.Join(".", Some(logger.Path, "logs").(string)))
		payload := &Payload{}
		if logger.Format == nil {
			logger.Format = defaultFormat
		}
		router.Use(func(c *gin.Context) {
			path := c.Request.URL.Path
			raw := c.Request.URL.RawQuery
			if raw != "" {
				path = path + "?" + raw
			}
			payload.Type = INT
			payload.StartUnix = time.Now().Unix()
			payload.IP = c.ClientIP()
			payload.Method = c.Request.Method
			payload.URL = path
			in := logger.Format(payload, c)
			journal.Info("%s", in)
			c.Next()
			payload.Type = OUT
			payload.EndUnix = time.Now().Unix()
			out := logger.Format(payload, c)
			journal.Info("%s", out)
		})
	}
}
