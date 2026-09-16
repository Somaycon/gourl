package model

import "gorm.io/gorm"

type Url struct {
	gorm.Model
	Code     string `json:"code" gorm:"uniqueIndex"`
	ShortUrl string `json:"short_url" binding:"required, url"`
	Url      string `json:"url"`
	Clicks   int    `gorm:"default:0" json:"clicks"`
}

type UrlRequest struct {
	Url string `json:"url"`
}
