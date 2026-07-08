package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT Claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// jwtSecret JWT 签名密钥，通过 InitJWTSecret 注入，避免硬编码
var jwtSecret = []byte("dev-only-jwt-secret-change-me")

// InitJWTSecret 初始化 JWT 签名密钥。应在应用启动时调用。
func InitJWTSecret(secret string) {
	if secret != "" {
		jwtSecret = []byte(secret)
	}
}

// GenerateToken 生成JWT Token
func GenerateToken(user User) (string, int64, error) {
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "kubeadm",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
