package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	port         *int
	ginMode      *string
	logToFile    *bool
	useGinLogger *bool
	timeFormat   = "2006-01-02T15:04:05GMT"
)

func init() {
	port = flag.Int("port", 8000, "Http running on the port")
	ginMode = flag.String("ginMode", "release", "Gin webframework running on release mode")
	useGinLogger = flag.Bool("useGinLogger", false, "Use gin logger instead of the one used in production")
}

func main() {
	flag.Parse()
	gin.SetMode(*ginMode)
	var route *gin.Engine
	if *useGinLogger {
		route = gin.Default()
	} else {
		route = gin.New()
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
