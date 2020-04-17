package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	CTEMERG   int = 0
	CTALERT   int = 1
	CTCRIT    int = 2
	CTERR     int = 3
	CTWARNING int = 4
	CTNOTICE  int = 5
	CTINFO    int = 6
	CTDEBUG   int = 7
)

type Ctlog struct {
	Fp       *os.File
	Loger    *log.Logger
	LogLevel int
}

func (c *Ctlog) SetLevel(level int) {
	c.LogLevel = level
}

func (c *Ctlog) Emerg(format string, args ...interface{}) {
	if c.LogLevel < CTEMERG {
		return
	}
	c.write("EMERG", format, args...)
}
func (c *Ctlog) Alert(format string, args ...interface{}) {
	if c.LogLevel < CTALERT {
		return
	}
	c.write("ALERT", format, args...)
}
func (c *Ctlog) Crit(format string, args ...interface{}) {
	if c.LogLevel < CTCRIT {
		return
	}
	c.write("CRIT", format, args...)
}
func (c *Ctlog) Error(format string, args ...interface{}) {
	if c.LogLevel < CTERR {
		return
	}
	c.write("ERROR", format, args...)
}
func (c *Ctlog) Warning(format string, args ...interface{}) {
	if c.LogLevel < CTWARNING {
		return
	}
	c.write("WARNING", format, args...)
}
func (c *Ctlog) Notice(format string, args ...interface{}) {
	if c.LogLevel < CTNOTICE {
		return
	}
	c.write("NOTICE", format, args...)
}
func (c *Ctlog) Info(format string, args ...interface{}) {
	if c.LogLevel < CTINFO {
		return
	}
	c.write("INFO", format, args...)
}
func (c *Ctlog) Debug(format string, args ...interface{}) {
	if c.LogLevel < CTDEBUG {
		return
	}
	c.write("DEBUG", format, args...)
}

func NewCTLog(filename string, useLogger bool) (*Ctlog, error) {
	// 创建一个文件并保存文件句柄
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Println("Create Log File Fail: ", err)
		return nil, nil
	}

	c := &Ctlog{}
	c.Fp = fp
	c.LogLevel = CTDEBUG

	if useLogger == true {
		c.Loger = log.New(c.Fp, "", log.Ldate|log.Ltime)
	}

	return c, nil
}

func (c *Ctlog) Close() {
	if c.Fp != nil {
		c.Fp.Close()
	}
}

func (c *Ctlog) write(level string, format string, args ...interface{}) {
	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(2)
	if ok {
		funcname = runtime.FuncForPC(pc).Name()      // main.(*myStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo

		filename = filepath.Base(filename) // /full/path/basename.go => basename.go
	}

	//log.Printf("%s:%d:%s: %s: %s\n", filename, line, funcname, level, fmt.Sprintf(format, args...))
	logStr := fmt.Sprintf("%s:%d:%s: [%s]: %s", filename, line, funcname, level, fmt.Sprintf(format, args...))

	if c.Loger != nil {
		c.Loger.Println(logStr)
		return
	}

	now := time.Now()
	thisTime := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	logLine := fmt.Sprintf("%s %s\n", thisTime, logStr)
	_, err := c.Fp.Write([]byte(logLine))
	if err != nil {
		fmt.Println("Write Log To LogFile Fail: ", err)
	}
}
