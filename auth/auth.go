package auth

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var authKey = "abc123456abc"

func GetToken(sp string) string {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["ID"] = "mujiburTest"
	token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	ts, err := token.SignedString([]byte(sp))
	if err != nil {
		log.Printf("Error to get signed string: %v\n", err)
		return ""
	}
	return ts
}

// Auth middleware returns handler that serves if the request comes in with valid JSON Web Token
//It will collect the authkey from the request and add that with GIN Context keys for future request
func Auth(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Keys = make(map[string]interface{})
		ctx.Keys["authKey"] = authKey
		_, err := jwt.ParseFromRequest(ctx.Request, func(token *jwt.Token) (interface{}, error) {
			key := ([]byte(secret))
			return key, nil
		})
		if err != nil {
			log.Println("Error: ", err)
			ctx.AbortWithError(401, err)
		}
	}
}
func VerifyAuthKey(key string) bool {
	if key != authKey {
		return false
	}
	return true
}
