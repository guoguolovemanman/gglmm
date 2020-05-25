package gglmm

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthKey 认证信息键类型
type AuthKey string

const (
	// AuthInfoRequestKey 认证信息请求建
	AuthInfoRequestKey AuthKey = "gglmm-auth-info"
)

// AuthInfo 认证信息
type AuthInfo struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}

// Authenticationable 可认证类型
type Authenticationable interface {
	AuthInfo() *AuthInfo
}

// GenerateAuthJWT 生成 Authorization Token
func GenerateAuthJWT(user Authenticationable, expires int64, secret string) (string, *jwt.StandardClaims, error) {
	jwtClaims, err := jwtGenerateClaims(user.AuthInfo(), expires)
	if err != nil {
		return "", jwtClaims, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", jwtClaims, err
	}
	return tokenString, jwtClaims, nil
}

// ParseAuthJWT 解析 Authorization Token
func ParseAuthJWT(tokenString string, secret string) (*AuthInfo, *jwt.StandardClaims, error) {
	jwtClaims, err := jwtParseClaims(tokenString, secret)
	if err != nil {
		return nil, nil, err
	}
	authInfo := &AuthInfo{}
	err = json.Unmarshal([]byte(jwtClaims.Subject), authInfo)
	if err != nil {
		return nil, nil, ErrAuthJWT
	}
	return authInfo, jwtClaims, nil
}

// GetAuthJWT 从请求里取 Authorization Token
func GetAuthJWT(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

// RequestWithAuthInfo 给请求设置 Authorization
func RequestWithAuthInfo(r *http.Request, authInfo *AuthInfo) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), AuthInfoRequestKey, authInfo))
}

// GetAuthInfo 从请求取 Authorization
func GetAuthInfo(r *http.Request) (*AuthInfo, error) {
	value := r.Context().Value(AuthInfoRequestKey)
	if value == nil {
		return nil, ErrAuthInfoNotFound
	}
	authInfo, ok := value.(*AuthInfo)
	if !ok {
		return nil, ErrAuthInfoNotFound
	}
	return authInfo, nil
}

// GetAuthType 从请求取 Authorization Type
func GetAuthType(r *http.Request) (string, error) {
	authInfo, err := GetAuthInfo(r)
	if err != nil {
		return "", err
	}
	return authInfo.Type, nil
}

// GetAuthID 从请求取 Authorization ID
func GetAuthID(r *http.Request, checkType string) (int64, error) {
	authInfo, err := GetAuthInfo(r)
	if err != nil {
		return 0, err
	}
	if authInfo.Type != checkType {
		return 0, ErrAuthType
	}
	return authInfo.ID, nil
}

const (
	// JWTExpires JWT失效时间
	JWTExpires int64 = 24 * 60 * 60
)

func jwtGenerateClaims(subject interface{}, expires int64) (*jwt.StandardClaims, error) {
	subjectBytes, err := json.Marshal(subject)
	if err != nil {
		return nil, err
	}
	jwtClaims := &jwt.StandardClaims{}
	now := time.Now().Unix()
	jwtClaims.IssuedAt = now
	jwtClaims.NotBefore = now
	jwtClaims.ExpiresAt = now + expires
	jwtClaims.Subject = string(subjectBytes)
	return jwtClaims, nil
}

func jwtParseClaims(tokenString string, secret string) (*jwt.StandardClaims, error) {
	jwtClaims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrAuthJWT
	}
	return jwtClaims, nil
}
