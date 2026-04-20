package handlers

import (
	"net/http"

	"github.com/unixadmin/anime/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if getSessionUserID(r) == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if getSessionUserID(r) == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if getSessionRole(r) != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}

}

func RequireAPIKey(queries *db.Queries) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			rawKey := authHeader[7:]
			accounts, err := queries.ListServiceAccounts(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			for _, a := range accounts {
				if bcrypt.CompareHashAndPassword([]byte(a.ApiKeyHash), []byte(rawKey)) == nil {
					next(w, r)
					return
				}
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}
