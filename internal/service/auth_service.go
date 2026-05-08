// Package service 业务逻辑层。
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/yinyin/myblog/internal/dto"
	"github.com/yinyin/myblog/internal/model"
	myjwt "github.com/yinyin/myblog/internal/pkg/jwt"
	"github.com/yinyin/myblog/internal/pkg/password"
	"github.com/yinyin/myblog/internal/repository"
)

// AuthService 鉴权服务接口
type AuthService interface {
	Login(ctx context.Context, req dto.LoginReq) (*dto.LoginResp, error)
}

type authService struct {
	users  repository.UserRepo
	signer *myjwt.Signer
}

// NewAuthService 构造
func NewAuthService(users repository.UserRepo, signer *myjwt.Signer) AuthService {
	return &authService{users: users, signer: signer}
}

// Login 用户名 + 密码登录
func (s *authService) Login(ctx context.Context, req dto.LoginReq) (*dto.LoginResp, error) {
	u, err := s.users.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// 不泄漏"用户存不存在",统一返回密码错误
			return nil, ErrInvalidPassword
		}
		return nil, fmt.Errorf("login: get user: %w", err)
	}

	if !u.IsActive() {
		return nil, ErrUserDisabled
	}
	if !password.Verify(u.PasswordHash, req.Password) {
		return nil, ErrInvalidPassword
	}

	access, refresh, exp, err := s.signer.Issue(u.ID, u.Username, u.Role)
	if err != nil {
		return nil, fmt.Errorf("login: issue token: %w", err)
	}

	return &dto.LoginResp{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresAt:    exp.Unix(),
		User:         toUserInfoResp(u),
	}, nil
}

func toUserInfoResp(u *model.User) dto.UserInfoResp {
	return dto.UserInfoResp{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Role:     u.Role,
	}
}
