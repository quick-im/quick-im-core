package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey = []byte("config.Config.Jwt.Key")

type MyCustomClaims struct {
	Sid uint32 `json:"sid"`
	jwt.RegisteredClaims
}

func ReleaseToken(sid uint32) (token string, err error) {
	// Create the Claims
	claims := MyCustomClaims{
		sid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "quick-im",
		},
	}

	tokenRet := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenRet.SignedString(mySigningKey)
	return
}

func ParseToken(token string) (*MyCustomClaims, error) {

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if tokenClaims != nil {
		if tokenClaims.Valid {
			fmt.Println("You look nice today")
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			fmt.Println("That's not even a token")
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			// Invalid signature
			fmt.Println("Invalid signature")
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*MyCustomClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err

}
