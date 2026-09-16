package repository

import (
	"context"
	"time"

	"github.com/Somaycon/gourl.git/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UrlRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewUrlRepository(db *gorm.DB, rdb *redis.Client) *UrlRepository {
	return &UrlRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *UrlRepository) Create(url *model.Url) error {
	return r.db.Create(url).Error
}

func (r *UrlRepository) FindByCode(code string) (*model.Url, error) {
	var url model.Url
	err := r.db.Where("code = ?", code).First(&url).Error

	return &url, err
}

func (r *UrlRepository) IncrementClicks(code string) {
	r.db.Model(&model.Url{}).Where("code = ?", code).UpdateColumn("clicks", gorm.Expr("clicks + ?", 1))
}

func (r *UrlRepository) GetCache(ctx context.Context, code string) (string, error) {
	return r.rdb.Get(ctx, code).Result()
}

func (r *UrlRepository) SetCache(ctx context.Context, code string, url string, duration time.Duration) error {
	return r.rdb.Set(ctx, code, url, duration).Err()
}
