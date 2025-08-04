package helpers

import (
	"erp/api/common"
	"erp/pkg/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtHelper interface {
	GenerateToken(claims common.Claims) (token string, err error)
	ExtractClaimsAdmin(tokenString string) (*common.Claims, error)
	GetToken(token string) string
}

type jwtHelper struct {
	secret []byte
}

func NewJwtHelper(
	config *config.AppConfig,
) JwtHelper{
	return &jwtHelper{
		secret:[]byte(config.Api.JwtSecret),
	}
}


func (s *jwtHelper) GenerateToken(claims common.Claims)(token string,err error){
	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	token, err = jwtToken.SignedString(s.secret)
	if err != nil {
		return 
	}
	return
}

func (s *jwtHelper)ExtractClaimsAdmin(tokenString string) (*common.Claims, error) {
	claims := &common.Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(tokenKey *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return claims, err
	}
	return claims, err
}

// func (s *jwtHelper)ExtractClaimsUserSession(tokenString string) (*common.UserSession, error) {
// 	claims := &common.UserSession{}
// 	_, err := jwt.ParseWithClaims(tokenString, claims, func(tokenKey *jwt.Token) (interface{}, error) {
// 		return s.frontendSecret, nil
// 	})
// 	if err != nil {
// 		return claims, err
// 	}
// 	return claims, err
// }

func (s *jwtHelper)GetToken(token string) string {
	return strings.TrimSpace(strings.Split(token, "Bearer")[1])
}
