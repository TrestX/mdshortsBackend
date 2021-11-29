package db

import (
	"MdShorts/pkg/entity"
)

type NewsService interface {
	GetNews(userId string) ([]entity.NewsDB, error)
	GetGlobalNews(country, language string) ([]entity.NewsDB, error)
	GetSearchNews(search string) ([]entity.NewsDB, error)
}
