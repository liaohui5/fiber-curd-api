package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string // token 类型

const (
	TokenTypeAccess  TokenType = "access"  // 访问token
	TokenTypeRefresh TokenType = "refresh" // 刷新token
)

type CustomClaims struct {
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// ConfigProvider 抽象配置来源, 测试时用 Mock 实现
type ConfigProvider interface {
	AccessTokenExpiredTime() int64
	AccessTokenSecret() string
	RefreshTokenExpiredTime() int64
	RefreshTokenSecret() string
}

// NowFunc 时间函数类型, 方便测试时注入固定时间
type NowFunc func() time.Time

// TokenIDFunc Token ID 生成函数类型, 方便测试时注入固定 ID
type TokenIDFunc func() string

// JWTService 封装所有 JWT 操作, 通过依赖注入解耦外部依赖
type JWTService struct {
	config    ConfigProvider
	nowFunc   NowFunc
	tokenIDFn TokenIDFunc
}

// NewJWTService 创建 JWTService，必须传入 ConfigProvider
// 可选传入 NowFunc 和 TokenIDFunc，未传则使用默认实现
type Option func(*JWTService)

func WithNowFunc(fn NowFunc) Option {
	return func(s *JWTService) {
		s.nowFunc = fn
	}
}

func WithTokenIDFunc(fn TokenIDFunc) Option {
	return func(s *JWTService) {
		s.tokenIDFn = fn
	}
}

func NewJWTService(config ConfigProvider, opts ...Option) *JWTService {
	svc := &JWTService{
		config:    config,
		nowFunc:   time.Now,
		tokenIDFn: func() string { return uuid.New().String() },
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// GenAccessToken 生成访问令牌(短期有效)
func (s *JWTService) GenAccessToken(uid string) (string, error) {
	now := s.nowFunc()
	expireDuration := time.Duration(s.config.AccessTokenExpiredTime()) * time.Minute

	claims := CustomClaims{
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   uid,
			ID:        s.tokenIDFn(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.AccessTokenSecret()))
}

// VerifyAccessToken 验证访问令牌
func (s *JWTService) VerifyAccessToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &CustomClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.AccessTokenSecret()), nil
		},
		jwt.WithTimeFunc(s.nowFunc),
	)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.TokenType != TokenTypeAccess {
			return "", errors.New("token type mismatch: expected access")
		}
		return claims.GetSubject()
	}

	return "", errors.New("invalid token")
}

// GenRefreshToken 生成刷新令牌（长期有效）
func (s *JWTService) GenRefreshToken(uid string) (string, error) {
	now := s.nowFunc()
	expireDuration := time.Duration(s.config.RefreshTokenExpiredTime()) * time.Hour

	claims := CustomClaims{
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   uid,
			ID:        s.tokenIDFn(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.RefreshTokenSecret()))
}

// VerifyRefreshToken 验证刷新令牌
func (s *JWTService) VerifyRefreshToken(refreshToken string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &CustomClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.RefreshTokenSecret()), nil
		},
		jwt.WithTimeFunc(s.nowFunc),
	)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.TokenType != TokenTypeRefresh {
			return "", errors.New("token type mismatch: expected refresh")
		}
		return claims.GetSubject()
	}

	return "", errors.New("invalid refresh token")
}

// RenewAccessToken 使用 Refresh Token 换取新的 Access Token
func (s *JWTService) RenewAccessToken(refreshToken string) (string, string, error) {
	uid, err := s.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := s.GenAccessToken(uid)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.GenRefreshToken(uid)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
