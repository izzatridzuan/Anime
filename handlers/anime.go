package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	cloudinaryHelper "github.com/unixadmin/anime/cloudinary"
	"github.com/unixadmin/anime/internal/db"
	"github.com/unixadmin/anime/models"
	"github.com/unixadmin/anime/templates"
	animetmpl "github.com/unixadmin/anime/templates/anime"
)

type AnimeHandler struct {
	queries *db.Queries
}

func NewAnimeHandler(queries *db.Queries) *AnimeHandler {
	return &AnimeHandler{queries: queries}
}

func (h *AnimeHandler) List(w http.ResponseWriter, r *http.Request) {
	animes, err := h.queries.ListAnime(r.Context())
	if err != nil {
		slog.Error("failed to list anime", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := make([]models.AnimeWithStudios, len(animes))
	for i, a := range animes {
		studios, err := h.queries.GetStudiosByAnime(r.Context(), a.ID)
		if err != nil {
			slog.Error("failed to get studios for anime", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result[i] = models.AnimeWithStudios{Anime: a, Studios: studios}
	}
	animetmpl.List(result, getSessionRole(r), getFlash(w, r)).Render(r.Context(), w)
}

func (h *AnimeHandler) New(w http.ResponseWriter, r *http.Request) {
	studios, err := h.queries.ListStudios(r.Context())
	if err != nil {
		slog.Error("failed to list studios for form", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	animetmpl.Form(db.Anime{}, false, studios, []int32{}, getSessionRole(r)).Render(r.Context(), w)
}

func (h *AnimeHandler) Create(w http.ResponseWriter, r *http.Request) {
	episodes, _ := strconv.Atoi(r.FormValue("episodes"))

	// Standard Input Validation
	if err := validateAnimeForm(r.FormValue("title"), r.FormValue("genre"), r.FormValue("status"), episodes); err != nil {
		templates.AlertError(err.Message).Render(r.Context(), w)
		return
	}

	// Date Validation
	releaseDate, dateErr := validateReleaseDate(r.FormValue("release_date"))
	if dateErr != nil {
		templates.AlertError(dateErr.Message).Render(r.Context(), w)
		return
	}

	//Image Validation
	file, header, err := r.FormFile("image")
	if err != nil {
		templates.AlertError("Image is required").Render(r.Context(), w)
		return
	}
	defer file.Close()

	if validErr := validateImageUpload(header); validErr != nil {
		templates.AlertError(validErr.Message).Render(r.Context(), w)
		return
	}

	imageUrl, err := cloudinaryHelper.UploadImage(r.Context(), file, header.Filename)
	if err != nil {
		templates.AlertError("Failed to upload image. Please try again.").Render(r.Context(), w)
		return
	}
	anime, err := h.queries.CreateAnime(r.Context(), db.CreateAnimeParams{
		Title:       r.FormValue("title"),
		Genre:       r.FormValue("genre"),
		Episodes:    int32(episodes),
		Status:      r.FormValue("status"),
		ImageUrl:    imageUrl,
		ReleaseDate: releaseDate,
	})
	if err != nil {
		templates.AlertError("Failed to create anime. Please try again").Render(r.Context(), w)
		return
	}

	for _, idStr := range r.Form["studio_ids"] {
		studioID, _ := strconv.Atoi(idStr)
		h.queries.AddStudioToAnime(r.Context(), db.AddStudioToAnimeParams{
			AnimeID:  anime.ID,
			StudioID: int32(studioID),
		})
	}
	logAudit(r.Context(), h.queries, getSessionUserID(r), "created", "anime", anime.ID)
	setFlash(w, r, "Anime created successfully")
	w.Header().Set("HX-Redirect", "/anime")
}

func (h *AnimeHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	anime, err := h.queries.GetAnime(r.Context(), int32(id))
	if err != nil {
		slog.Error("anime not found", "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	studios, err := h.queries.ListStudios(r.Context())
	if err != nil {
		slog.Error("failed to list studios for edit", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	assigned, err := h.queries.GetStudiosByAnime(r.Context(), int32(id))
	if err != nil {
		slog.Error("failed to get assigned studios", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	assignedIDs := make([]int32, len(assigned))
	for i, s := range assigned {
		assignedIDs[i] = s.ID
	}
	animetmpl.Form(anime, true, studios, assignedIDs, getSessionRole(r)).Render(r.Context(), w)
}

func (h *AnimeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	episodes, _ := strconv.Atoi(r.FormValue("episodes"))

	if err := validateAnimeForm(r.FormValue("title"), r.FormValue("genre"), r.FormValue("status"), episodes); err != nil {
		templates.AlertError(err.Message).Render(r.Context(), w)
		return
	}
	releaseDate, dateErr := validateReleaseDate(r.FormValue("release_date"))
	if dateErr != nil {
		templates.AlertError(dateErr.Message).Render(r.Context(), w)
		return
	}
	existing, err := h.queries.GetAnime(r.Context(), int32(id))
	if err != nil {
		templates.AlertError("Anime not found.").Render(r.Context(), w)
		return
	}

	imageUrl := existing.ImageUrl
	file, header, err := r.FormFile("image")
	if err == nil {
		// new file uploaded — validate and upload
		defer file.Close()
		if validErr := validateImageUpload(header); validErr != nil {
			templates.AlertError(validErr.Message).Render(r.Context(), w)
			return
		}
		if existing.ImageUrl != "" {
			cloudinaryHelper.DeleteImage(r.Context(), existing.ImageUrl)
		}
		imageUrl, err = cloudinaryHelper.UploadImage(r.Context(), file, header.Filename)
		if err != nil {
			templates.AlertError("Failed to upload image. Please try again.").Render(r.Context(), w)
			return
		}
	} else if existing.ImageUrl == "" {
		// no file uploaded and no existing image — required
		templates.AlertError("Image is required").Render(r.Context(), w)
		return
	}

	_, err = h.queries.UpdateAnime(r.Context(), db.UpdateAnimeParams{
		ID:          int32(id),
		Title:       r.FormValue("title"),
		Genre:       r.FormValue("genre"),
		Episodes:    int32(episodes),
		Status:      r.FormValue("status"),
		ImageUrl:    imageUrl,
		ReleaseDate: releaseDate,
	})
	if err != nil {
		templates.AlertError("Failed to update anime. Please try again.").Render(r.Context(), w)
		return
	}
	h.queries.DeleteAnimeStudios(r.Context(), int32(id))
	for _, idStr := range r.Form["studio_ids"] {
		studioID, _ := strconv.Atoi(idStr)
		h.queries.AddStudioToAnime(r.Context(), db.AddStudioToAnimeParams{
			AnimeID:  int32(id),
			StudioID: int32(studioID),
		})
	}
	logAudit(r.Context(), h.queries, getSessionUserID(r), "updated", "anime", int32(id))
	setFlash(w, r, "Anime updated successfully")
	w.Header().Set("HX-Redirect", "/anime")
}

func (h *AnimeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	err := h.queries.ArchiveAnime(r.Context(), int32(id))
	if err != nil {
		slog.Error("failed to archive anime", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logAudit(r.Context(), h.queries, getSessionUserID(r), "archived", "anime", int32(id))
	setFlash(w, r, "Anime deleted successfully")
	w.Header().Set("HX-Redirect", "/anime")
}
