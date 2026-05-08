package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/yinyin/myblog/internal/dto"
	"github.com/yinyin/myblog/internal/model"
	"github.com/yinyin/myblog/internal/repository"
)

// ProfileService 个人资料服务
type ProfileService interface {
	Get(ctx context.Context) (*dto.ProfileResp, error)
	Update(ctx context.Context, req dto.ProfileUpdateReq) (*dto.ProfileResp, error)
}

type profileService struct {
	repo repository.ProfileRepo
}

// NewProfileService 构造
func NewProfileService(repo repository.ProfileRepo) ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) Get(ctx context.Context) (*dto.ProfileResp, error) {
	p, err := s.repo.Get(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// 首次访问时返回空 profile(前端可以展示默认值)
			return &dto.ProfileResp{}, nil
		}
		return nil, fmt.Errorf("profile get: %w", err)
	}
	return toProfileResp(p), nil
}

func (s *profileService) Update(ctx context.Context, req dto.ProfileUpdateReq) (*dto.ProfileResp, error) {
	p := &model.Profile{
		Name:     req.Name,
		Bio:      req.Bio,
		Avatar:   req.Avatar,
		Email:    req.Email,
		GitHub:   req.GitHub,
		Twitter:  req.Twitter,
		LinkedIn: req.LinkedIn,
		Website:  req.Website,
	}
	if err := s.repo.Upsert(ctx, p); err != nil {
		return nil, fmt.Errorf("profile update: %w", err)
	}
	return toProfileResp(p), nil
}

func toProfileResp(p *model.Profile) *dto.ProfileResp {
	return &dto.ProfileResp{
		Name:     p.Name,
		Bio:      p.Bio,
		Avatar:   p.Avatar,
		Email:    p.Email,
		GitHub:   p.GitHub,
		Twitter:  p.Twitter,
		LinkedIn: p.LinkedIn,
		Website:  p.Website,
	}
}
