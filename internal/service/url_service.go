package service

import (
	"context"
	"errors"
	"time"

	"github.com/Somaycon/gourl.git/internal/model"
	"github.com/Somaycon/gourl.git/internal/pkg/generator"
	"github.com/Somaycon/gourl.git/internal/repository"
)

type UrlService struct {
	repo    *repository.UrlRepository
	baseUrl string
}

func NewUrlService(repo *repository.UrlRepository, baseUrl string) *UrlService {
	return &UrlService{
		repo:    repo,
		baseUrl: baseUrl,
	}
}

func (s *UrlService) Shorten(ctx context.Context, originalUrl string) (*model.Url, error) {
	var code string
	for {
		generated, err := generator.GenerateCode(6)
		if err != nil {
			return nil, err
		}
		_, err = s.repo.FindByCode(generated)

		if err != nil {
			code = generated
			break
		}
	}
	newUrl := &model.Url{
		Code:     code,
		ShortUrl: s.baseUrl + code,
		Url:      originalUrl,
	}
	if err := s.repo.Create(newUrl); err != nil {
		return nil, err
	}

	return newUrl, nil
}

func (s *UrlService) GetOrininalUrl(ctx context.Context, code string) (string, error) {
	cachedUrl, err := s.repo.GetCache(ctx, code)
	if err == nil {
		go s.repo.GetCache(ctx, code)
		s.repo.IncrementClicks(code)
		return cachedUrl, nil
	}

	urlModel, err := s.repo.FindByCode(code)
	if err != nil {
		return "", errors.New("Url not found")
	}
	s.repo.IncrementClicks(code)
	_ = s.repo.SetCache(ctx, code, urlModel.Url, 24*time.Hour)
	return urlModel.Url, nil
}
