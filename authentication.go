package gglmm

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthorizationKey 认证信息键类型
type AuthorizationKey string

const (
	// AuthorizationRequestKey 认证信息请求建
	AuthorizationRequestKey AuthorizationKey = "gglmm-authorization"
)

// Authorization 认证信息
type Authorization struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}

// Authenticationable 可认证类型
type Authenticationable interface {
	Authorization() *Authorization
}

// GenerateAuthorizationToken 生成 Authorization Token
func GenerateAuthorizationToken(user Authenticationable, expires int64, secret string) (string, *jwt.StandardClaims, error) {
	jwtClaims, err := jwtGenerateClaims(user.Authorization(), expires)
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

// ParseAuthorizationToken 解析 Authorization Token
func ParseAuthorizationToken(tokenString string, secret string) (*Authorization, *jwt.StandardClaims, error) {
	jwtClaims, err := jwtParseClaims(tokenString, secret)
	if err != nil {
		return nil, nil, err
	}
	authorization := &Authorization{}
	err = json.Unmarshal([]byte(jwtClaims.Subject), authorization)
	if err != nil {
		return nil, nil, errors.New("cannot convert subject to jwtSubject")
	}
	return authorization, jwtClaims, nil
}

// GetAuthorizationToken 从请求里取 Authorization Token
func GetAuthorizationToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

// RequestWithAuthorization 给请求设置 Authorization
func RequestWithAuthorization(r *http.Request, authorization *Authorization) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), AuthorizationRequestKey, authorization))
}

// GetAuthorization 从请求取 Authorization
func GetAuthorization(r *http.Request) (*Authorization, error) {
	value := r.Context().Value(AuthorizationRequestKey)
	if value == nil {
		return nil, errors.New("jwtsubject error")
	}
	authorization, ok := value.(*Authorization)
	if !ok {
		return nil, errors.New("jwtsubject type error")
	}
	return authorization, nil
}

// GetAuthorizationType 从请求取 Authorization Type
func GetAuthorizationType(r *http.Request) (string, error) {
	authorization, err := GetAuthorization(r)
	if err != nil {
		return "", err
	}
	return authorization.Type, nil
}

// GetAuthorizationID 从请求取 Authorization ID
func GetAuthorizationID(r *http.Request, checkType string) (int64, error) {
	authorization, err := GetAuthorization(r)
	if err != nil {
		return 0, err
	}
	if authorization.Type != checkType {
		return 0, errors.New("jwtsubject type check fail")
	}
	return authorization.ID, nil
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
		return nil, errors.New("cannot convert claim to jwt.StandardClaims")
	}
	return jwtClaims, nil
}
