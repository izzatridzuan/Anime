package models

import (
	"time"

	db "github.com/unixadmin/anime/internal/db"
)

type AnimeWithStudios struct {
	Anime   db.Anime
	Studios []db.Studio
}

type StudioResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type AnimeResponse struct {
	ID          int32            `json:"id"`
	Title       string           `json:"title"`
	Genre       string           `json:"genre"`
	Episodes    int32            `json:"episodes"`
	Status      string           `json:"status"`
	ImageUrl    string           `json:"image_url"`
	ReleaseDate *time.Time       `json:"release_date"`
	CreatedAt   time.Time        `json:"created_at"`
	Studios     []StudioResponse `json:"studios"`
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
