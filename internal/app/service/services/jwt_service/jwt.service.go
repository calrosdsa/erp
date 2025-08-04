package jwt_service

import (
	"erp/api/common"
	"erp/internal/app/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secret []byte
	frontendSecret []byte
}


func NewJwtService(configService *config.ConfigService) *JwtService{
	secret := configService.GetApiOptions().JwtSecret
	return &JwtService{
		secret: []byte(secret),
	}
}

func (s *JwtService) GenerateToken(claims common.Claims)(token string,err error){
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	token, err = jwtToken.SignedString(s.secret)
	if err != nil {
		return token,err
	}
	return
}

func (s *JwtService)ExtractClaimsAdmin(tokenString string) (*common.Claims, error) {
	claims := &common.Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(tokenKey *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return claims, err
	}
	return claims, err
}

// func (s *JwtService)ExtractClaimsUserSession(tokenString string) (*common.UserSession, error) {
// 	claims := &common.UserSession{}
// 	_, err := jwt.ParseWithClaims(tokenString, claims, func(tokenKey *jwt.Token) (interface{}, error) {
// 		return s.frontendSecret, nil
// 	})
// 	if err != nil {
// 		return claims, err
// 	}
// 	return claims, err
// }

func (s *JwtService)GetToken(token string) string {
	return strings.TrimSpace(strings.Split(token, "Bearer")[1])
}
