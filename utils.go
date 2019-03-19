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
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Some get or a default value
func Some(target interface{}, initValue interface{}) interface{} {
	if target != nil && target != "" && target != 0 {
		return target
	}
	return initValue
}

// LeftV -
func LeftV(left interface{}, right interface{}) interface{} {
	return left
}

// createLog -
func createLog(path string) io.Writer {
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0600)
	writer := io.MultiWriter(f, os.Stdout)
	return writer
}

// getLogFile -
func getLogFile(level LOGLEVEL, basePath string) string {
	var filePath string
	levelStr := levelStr(level)
	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if filePath != "" {
			return nil
		}
		fileName := info.Name()
		fileSize := info.Size()
		levelMatch := strings.HasPrefix(fileName, levelStr)
		sizeMatch := false
		if level == SYSLEVEL {
			sizeMatch = fileSize < SYSSTROBE
		} else if level == USERLEVEL {
			sizeMatch = fileSize < USERSTROBE
		}

		if levelMatch && sizeMatch {
			filePath = path
		}
		return nil
	})
	if filePath != "" {
		return filePath
	}
	// create level log file
	fileName := time.Now().Format("2006.01.02 15.04.05")
	fileName = fmt.Sprintf("%s_"+fileName+".log", levelStr)
	filePath = path.Join(basePath, fileName)
	return filePath
}

// levelStr -
func levelStr(level LOGLEVEL) string {
	var levelStr string
	switch level {
	case SYSLEVEL:
		levelStr = "0"
	case USERLEVEL:
		levelStr = "1"
	}
	return levelStr
}
