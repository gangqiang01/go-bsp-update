package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
	jwt.StandardClaims
}

var jwtSecret = []byte("apphub_edge")
var expiredTime = 8 * time.Hour

// 生成jwt的token
func MakeCliamsToken(o UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, o)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

// 解密jwt的token
func ParseCliamsToken(token string) (*UserClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*UserClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// 生成token,并返回token,过期时间1小时
func GenerateToken(name string) (string, error) {
	claims := UserClaims{
		Name:       name,
		CreateTime: time.Now(),
		StandardClaims: jwt.StandardClaims{
			// 过期时间8小时
			ExpiresAt: time.Now().Add(expiredTime).Unix(),
		},
	}
	return MakeCliamsToken(claims)
}
