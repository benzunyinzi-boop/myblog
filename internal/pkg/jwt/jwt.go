// Package jwt 签发与校验 access / refresh token。
package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 错误集合(service 层按需映射到 errcode)
var (
	ErrInvalidToken = errors.New("jwt: invalid token")
	ErrExpiredToken = errors.New("jwt: token expired")
	ErrWrongType    = errors.New("jwt: wrong token type")
)

// 令牌类型
const (
	TypeAccess  = "access"
	TypeRefresh = "refresh"
)

// Claims 自定义载荷
type Claims struct {
	UserID   uint64 `json:"uid"`
	Username string `json:"usr"`
	Role     string `json:"role"`
	Type     string `json:"typ"` // access / refresh
	jwt.RegisteredClaims
}

// Signer 负责签发与校验。
type Signer struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

// Config Signer 配置(对应 config.JWTConfig)
type Config struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string
}

// NewSigner 构造
func NewSigner(cfg Config) (*Signer, error) {
	if len(cfg.Secret) < 8 {
		return nil, fmt.Errorf("jwt secret too short (>=8 required)")
	}
	if cfg.AccessTTL <= 0 {
		cfg.AccessTTL = 2 * time.Hour
	}
	if cfg.RefreshTTL <= 0 {
		cfg.RefreshTTL = 168 * time.Hour
	}
	if cfg.Issuer == "" {
		cfg.Issuer = "myblog"
	}
	return &Signer{
		secret:     []byte(cfg.Secret),
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
		issuer:     cfg.Issuer,
	}, nil
}

// Issue 签发 access + refresh 一对 token。
// 返回 accessToken, refreshToken, accessExpiresAt
func (s *Signer) Issue(userID uint64, username, role string) (string, string, time.Time, error) {
	now := time.Now()
	accessExp := now.Add(s.accessTTL)

	access, err := s.sign(userID, username, role, TypeAccess, now, accessExp)
	if err != nil {
		return "", "", time.Time{}, err
	}
	refresh, err := s.sign(userID, username, role, TypeRefresh, now, now.Add(s.refreshTTL))
	if err != nil {
		return "", "", time.Time{}, err
	}
	return access, refresh, accessExp, nil
}

// Verify 校验 token,返回 claims。
// wantType 传空表示不强制类型;常用:Verify(..., TypeAccess) 校验访问令牌。
func (s *Signer) Verify(tokenStr, wantType string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	}, jwt.WithIssuer(s.issuer))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, ErrInvalidToken
	}
	if wantType != "" && claims.Type != wantType {
		return nil, ErrWrongType
	}
	return claims, nil
}

func (s *Signer) sign(uid uint64, usr, role, typ string, now, exp time.Time) (string, error) {
	claims := &Claims{
		UserID:   uid,
		Username: usr,
		Role:     role,
		Type:     typ,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", uid),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(s.secret)
}
