package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	port         *int
	ginMode      *string
	logToFile    *bool
	useGinLogger *bool
	logFile      = "logs/web-server.log"
	timeFormat   = "2006-01-02T15:04:05GMT"
)

func init() {
	port = flag.Int("port", 8000, "Http running on the port")
	ginMode = flag.String("ginMode", "release", "Gin webframework running on release mode")
	logToFile = flag.Bool("logToFile", true, "Log write to file")
	useGinLogger = flag.Bool("useGinLogger", false, "Use gin logger instead of the one used in production")
}

func logger(c *gin.Context) {
	var start time.Time
	const logFormat = "%s " + // Timestamp
		"%s " + // Client ip
		"%d " + // Response code
		"%v " + // Response Duration
		`"%s %s %s" ` + // Request method, path and protocol
		"%d " // Response size

	defer func() {
		const INTERNAL_SERVER_ERROR = 500
		if err := recover(); err != nil {
			duration := time.Now().Sub(start)
			log.Printf(logFormat+"\n%v\n%s",
				start.Format(timeFormat),
				c.ClientIP(),
				INTERNAL_SERVER_ERROR,
				duration,
				c.Request.Method,
				c.Request.URL.Path,
				c.Request.Proto,
				c.Writer.Size(),
				err,
				debug.Stack(),
			)
			c.AbortWithStatus(INTERNAL_SERVER_ERROR)
		}
	}()
	start = time.Now()
	c.Next()
	duration := time.Now().Sub(start)
	log.Printf(logFormat+"\n%s",
		start.Format(timeFormat),
		c.ClientIP(),
		c.Writer.Status(),
		duration,
		c.Request.Method,
		c.Request.URL.Path,
		c.Request.Proto,
		c.Writer.Size(),
		c.Errors.String(),
	)
}

func main() {
	flag.Parse()
	if *logToFile {
		logWriteTo, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("error opening to log file: %v", err)
		}
		log.SetOutput(logWriteTo)
	}
	gin.SetMode(*ginMode)
	var route *gin.Engine
	if *useGinLogger {
		route = gin.Default()
	} else {
		route = gin.New()
		route.Use(logger)
	}

	route.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Its an entry point :)"})
	})

	if err := route.Run(fmt.Sprintf(":%d", *port)); err != nil {
		logFatal("Http not running: ", err)
	} else {
		logPrintln("shutting down")
	}
}
func logTime() {
	fmt.Fprintf(os.Stderr, "%s ", time.Now().Format(timeFormat))
}
func logPrintln(args ...interface{}) {
	logTime()
	log.Println(args...)
}
func logFatal(args ...interface{}) {
	logTime()
	fmt.Fprintf(os.Stderr, "Error: ")
	log.Fatal(args...)
}
func logFatalf(format string, args ...interface{}) {
	logTime()
	fmt.Fprintf(os.Stderr, "Error: ")
	log.Fatalf(fmt.Sprintf(format, args...))
}
