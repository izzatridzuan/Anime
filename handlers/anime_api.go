package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/unixadmin/anime/internal/db"
	"github.com/unixadmin/anime/models"
)

func getPaginationParams(r *http.Request) (page, pageSize, offset int) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ = strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset = (page - 1) * pageSize
	return
}

// APIList godoc
// @Summary      List all anime
// @Description  Returns a paginated list of anime with optional filters
// @Tags         anime
// @Produce      json
// @Param        title      query     string  false  "Filter by title"
// @Param        genre      query     string  false  "Filter by genre"
// @Param        status     query     string  false  "Filter by status (ongoing, completed, upcoming)"
// @Param        season     query     string  false  "Filter by season (spring, summer, autumn, winter)"
// @Param        page       query     int     false  "Page number (default 1)"
// @Param        page_size  query     int     false  "Page size (default 10)"
// @Success      200  {object}  models.PaginatedResponse[models.AnimeResponse]
// @Failure      500  {object}  map[string]string
// @Router       /api/anime [get]
func (h *AnimeHandler) APIList(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	genre := r.URL.Query().Get("genre")
	status := r.URL.Query().Get("status")
	season := r.URL.Query().Get("season")

	startDate, endDate, seasonErr := getSeasonDates(season)

	if seasonErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "` + seasonErr.Message + `"}`))
		return
	}
	page, pageSize, offset := getPaginationParams(r)

	total, err := h.queries.CountFilteredAnime(r.Context(), db.CountFilteredAnimeParams{
		Column1: title,
		Column2: genre,
		Column3: status,
		Column4: startDate,
		Column5: endDate,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "failed to count anime"}`))
		return
	}

	animes, err := h.queries.FilterAnimePaginated(r.Context(), db.FilterAnimePaginatedParams{
		Column1: title,
		Column2: genre,
		Column3: status,
		Column4: startDate,
		Column5: endDate,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "failed to fetch anime"}`))
		return
	}

	result := make([]models.AnimeResponse, len(animes))
	for i, a := range animes {
		studios, err := h.queries.GetStudiosByAnime(r.Context(), a.ID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "failed to fetch studios"}`))
			return
		}
		studioResp := make([]models.StudioResponse, len(studios))
		for j, s := range studios {
			studioResp[j] = models.StudioResponse{ID: s.ID, Name: s.Name}
		}
		var releaseDate *time.Time
		if a.ReleaseDate.Valid {
			releaseDate = &a.ReleaseDate.Time
		}
		result[i] = models.AnimeResponse{
			ID:          a.ID,
			Title:       a.Title,
			Genre:       a.Genre,
			Episodes:    a.Episodes,
			Status:      a.Status,
			ImageUrl:    a.ImageUrl,
			ReleaseDate: releaseDate,
			CreatedAt:   a.CreatedAt.Time,
			Studios:     studioResp,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.PaginatedResponse[models.AnimeResponse]{
		Data:       result,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
	})
}
