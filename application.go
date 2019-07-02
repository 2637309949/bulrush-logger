// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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

func createLogger(path string) *journal.Journal {
	journal := journal.CreateLogger(
		journal.HTTP,
		nil,
		[]*journal.Transport{
			&journal.Transport{
				Dirname: path,
				Level:   journal.HTTP,
				Maxsize: journal.Maxsize,
			},
			&journal.Transport{
				Level: journal.HTTP,
			},
		},
	)
	return journal
}

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
func (logger *Logger) Plugin(cfg *bulrush.Config, router *gin.RouterGroup) {
	logger.Path = Some(logger.Path, cfg.Log.Path).(string)
	journal := createLogger(path.Join(Some(logger.Path, "logs").(string), "http"))
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
		journal.HTTP("%s", in)
		c.Next()
		payload.Type = OUT
		payload.EndUnix = time.Now().Unix()
		out := logger.Format(payload, c)
		journal.HTTP("%s", out)
	})
}
