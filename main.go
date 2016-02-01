package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golangTest/auth"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	port           *int
	ginMode        *string
	logToFile      *bool
	useGinLogger   *bool
	logFile        = "logs/web-server.log"
	timeFormat     = "2006-01-02T15:04:05GMT"
	secretpassword = "holabolaNotLikeaBull"
)

func init() {
	port = flag.Int("port", 8000, "Http running on the port")
	ginMode = flag.String("ginMode", "release", "Gin webframework running on release mode by default")
	logToFile = flag.Bool("logToFile", true, "Log write to file")
	useGinLogger = flag.Bool("useGinLogger", false, "Use gin logger instead of the one used in production")
}

//logger is the middleware to run with every request
//It will log all the request information like timestamp, clientIp, response code,
// response duration, method, path, protocol and response size
//It will also help to recover any sort panic, fatal etc
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
		if logWriteTo != nil {
			log.SetOutput(logWriteTo)
		}
	}
	gin.SetMode(*ginMode)
	var route *gin.Engine
	if *useGinLogger {
		route = gin.Default()
	} else {
		route = gin.New()
		route.Use(logger)
	}

	route.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Its an entry point :)"})
	})

	route.GET("/token", func(ctx *gin.Context) {
		tokenString := auth.GetToken(secretpassword)
		ctx.JSON(200, gin.H{"token": tokenString})
	})
	route.GET("/dbaccessor", func(ctx *gin.Context) {
		ctx.JSON(200, getFakeDbData())
	})
	route.Use(auth.Auth(secretpassword))
	route.POST("/auth", func(ctx *gin.Context) {
		ak := ctx.Request.FormValue("authkey")
		if ak == "" {
			ctx.JSON(401, "No auth key")
		} else if !auth.VerifyAuthKey(ak) {
			ctx.JSON(401, "Wrong key")
		} else {
			ctx.Redirect(http.StatusFound, "/user")
		}
	})
	route.GET("/user", func(ctx *gin.Context) {
		key := ctx.MustGet("authKey")
		udetails := dbAccessor(key.(string))
		ctx.JSON(200, gin.H{"user": udetails})
	})
	route.GET("/user/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		ctx.JSON(200, find(id))
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

// dbAccessor returns the user profile if the authkey matched with the
//request.
func dbAccessor(key string) interface{} {
	str := []byte(`{"abc123456abc":{"username": "Mujibur", "id": 9001, "country": "malaysia", "mobile": "+019378646"}}`)
	var userdata map[string]interface{}
	err := json.Unmarshal(str, &userdata)
	if err != nil {
		log.Println("Error on unmarshalling:", err)
	}
	return userdata[key]
}

//User is the important struct which will be storing the content from JSON
//After unmarshalling the JSON data it will store to User struct
type User struct {
	Id      int    `json:"id"`
	Name    string `json:"username"`
	Country string `json:"country"`
	Mobile  string `json:"mobile"`
}

// find function used to find the user information from a json data
// It will return the data which will match with request /user/:id
func find(idStr string) *User {
	str := []byte(`[{"username": "Mujibur", "id": 9001, "country": "malaysia", "mobile": "+019378646"},{"username": "Luis", "id": 9002, "country": "Spain", "mobile": "+1212121"},{"username": "Holabola", "id": 9003, "country": "HongPong", "mobile": "+1234"},{"username": "XXX", "id": 9004, "country": "SSS", "mobile": "+4234234"}]`)
	var UserList []*User
	err := json.Unmarshal(str, &UserList)
	if err != nil {
		log.Println("Error on unmarshalling:", err)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error to convert: ", err)
	}
	for _, u := range UserList {
		if u.Id == id {
			return u
		}
	}
	return nil
}

//getFakeDbData return the db content which has hard coded to db.config
func getFakeDbData() interface{} {
	jsonFile, err := ioutil.ReadFile("./db.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		return nil
	}
	var dbConfig interface{}
	err = json.Unmarshal(jsonFile, &dbConfig)
	if err != nil {
		log.Println("Error on unmarshalling: ", err)
	}
	return dbConfig
}
