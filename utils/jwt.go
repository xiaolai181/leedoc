package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("lennon")

//Gentoken 生成token
func Gentoken(username string) (string, error) {
	claims := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(MySecret)
}

//parseToken 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

const TokenExpireDuration_xsrf = time.Minute * 5

var MySecret_xsrf = []byte("lennon_xsrf")

type IPClaims struct {
	IPName string `json:"ip"`
	jwt.StandardClaims
}

//生成
func GenXsrf(ip string) (string, error) {
	claims := MyClaims{
		ip,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration_xsrf).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(MySecret_xsrf)
}

//parsexsrf 解析Xsrf
func ParseXsrf(tokenString string) (*IPClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret_xsrf, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*IPClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
