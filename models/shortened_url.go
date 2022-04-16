package models

import (
	"time"
)

type ShortenedUrl struct {
	ID          string `gorm:"primarykey"`
	OriginalUrl string
	CreatedAt   time.Time
}

type Click struct {
	ID             uint `gorm:"primarykey"`
	ShortenedUrlId string
	ShortenedUrl   ShortenedUrl `gorm:"foreignKey:ShortenedUrlId"`
	CreatedAt      time.Time
	RawRequest     string
}
