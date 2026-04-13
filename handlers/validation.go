package handlers

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const maxImageSize = 1048576 // 1MB in bytes

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func validateAnimeForm(title, genre, status string, episodes int) *ValidationError {
	if len(title) == 0 || len(title) > 200 {
		return &ValidationError{Field: "title", Message: "Title must be between 1 and 200 characters."}

	}
	if len(genre) == 0 || len(genre) > 100 {
		return &ValidationError{Field: "genre", Message: "Genre must be between 1 and 100 characters."}

	}
	if episodes < 0 || episodes > 10000 {
		return &ValidationError{Field: "episodes", Message: "Episodes must be between 0 and 10000"}
	}
	validStatuses := map[string]bool{"ongoing": true, "completed": true, "upcoming": true}
	if !validStatuses[status] {
		return &ValidationError{Field: "status", Message: "Status must be ongoing, completed or upcoming"}
	}
	return nil
}

func validateStudioForm(name string) *ValidationError {
	if len(name) == 0 || len(name) > 200 {
		return &ValidationError{Field: "name", Message: "Studio name must be between 1 and 200 characters"}
	}
	return nil
}

func validateImageUpload(header *multipart.FileHeader) *ValidationError {
	if header == nil {
		return &ValidationError{Field: "image", Message: "Image is required"}
	}
	if header.Size > maxImageSize {
		return &ValidationError{Field: "image", Message: "Image must be less than 1MB"}
	}
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}
	if !allowed[header.Header.Get("Content-Type")] {
		return &ValidationError{Field: "image", Message: "Image must be JPEG, PNG or GIF"}
	}
	return nil
}

func validateReleaseDate(dateStr string) (pgtype.Timestamp, *ValidationError) {
	if dateStr == "" {
		return pgtype.Timestamp{}, &ValidationError{Field: "release_date", Message: "Release date is required"}
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return pgtype.Timestamp{}, &ValidationError{Field: "release_date", Message: "Release date is invalid"}
	}
	return pgtype.Timestamp{Time: t, Valid: true}, nil
}

func getSeasonDates(season string) (pgtype.Timestamp, pgtype.Timestamp, *ValidationError) {
	if season == "" {
		return pgtype.Timestamp{Valid: false}, pgtype.Timestamp{Valid: false}, nil
	}

	year := time.Now().Year()
	var start, end time.Time

	switch season {
	case "spring":
		start = time.Date(year, time.March, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.May, 31, 23, 59, 59, 0, time.UTC)
	case "summer":
		start = time.Date(year, time.June, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.August, 31, 23, 59, 59, 0, time.UTC)
	case "autumn":
		start = time.Date(year, time.September, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.November, 30, 23, 59, 59, 0, time.UTC)
	case "winter":
		start = time.Date(year, time.December, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(year+1, time.February, 28, 23, 59, 59, 0, time.UTC)
	default:
		return pgtype.Timestamp{Valid: false}, pgtype.Timestamp{Valid: false}, &ValidationError{
			Field:   "season",
			Message: "Season must be spring, summer, autumn or winter",
		}
	}

	return pgtype.Timestamp{Time: start, Valid: true},
		pgtype.Timestamp{Time: end, Valid: true},
		nil
}

func validatePasswordComplexity(password string) *ValidationError {
	if len(password) < 12 {
		return &ValidationError{Field: "password", Message: "Password must be at least 12 characters"}
	}

	var hasUpper, hasLower, hasDigit, hasSymbol bool
	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		default:
			hasSymbol = true
		}
	}

	if !hasUpper {
		return &ValidationError{Field: "password", Message: "Password must contain at least one uppercase letter"}
	}
	if !hasLower {
		return &ValidationError{Field: "password", Message: "Password must contain at least one lowercase letter"}
	}
	if !hasDigit {
		return &ValidationError{Field: "password", Message: "Password must contain at least one digit"}
	}
	if !hasSymbol {
		return &ValidationError{Field: "password", Message: "Password must contain at least one symbol"}
	}
	return nil
}

func validateUserForm(name, email, role string) *ValidationError {
	if len(name) == 0 || len(name) > 200 {
		return &ValidationError{Field: "name", Message: "Name must be between 1 and 200 characters"}
	}
	if len(email) == 0 || len(email) > 200 {
		return &ValidationError{Field: "email", Message: "Email is required"}
	}
	validRoles := map[string]bool{"admin": true, "user": true}
	if !validRoles[role] {
		return &ValidationError{Field: "role", Message: "Role must be admin or user"}
	}
	return nil
}
