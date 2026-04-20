package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/unixadmin/anime/docs"
	"github.com/unixadmin/anime/handlers"
	"github.com/unixadmin/anime/internal/db"
	"github.com/unixadmin/anime/logger"
)

// @title           Anime Backoffice API
// @version         1.0
// @description     API for managing anime and studios
// @host            localhost:8080
// @BasePath        /

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := logger.Init(); err != nil {
		log.Fatal("failed to init logger:", err)
	}
	handlers.InitSessionStore()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	animeHandler := handlers.NewAnimeHandler(queries)
	studioHandler := handlers.NewStudioHandler(queries)
	authHandler := handlers.NewAuthHandler(queries)
	userHandler := handlers.NewUserHandler(queries)
	profileHandler := handlers.NewProfileHandler(queries)
	auditLogHandler := handlers.NewAuditLogHandler(queries)
	serviceAccountHandler := handlers.NewServiceAccountHandler(queries)
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("GET /login", authHandler.LoginPage)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("GET /logout", authHandler.Logout)
	mux.HandleFunc("GET /change-password", handlers.RequireLogin(authHandler.ChangePasswordPage))
	mux.HandleFunc("POST /change-password", handlers.RequireLogin(authHandler.ChangePassword))

	// Anime routes
	mux.HandleFunc("GET /anime", handlers.RequireLogin(animeHandler.List))
	mux.HandleFunc("GET /anime/new", handlers.RequireLogin(animeHandler.New))
	mux.HandleFunc("POST /anime", handlers.RequireLogin(animeHandler.Create))
	mux.HandleFunc("GET /anime/{id}/edit", handlers.RequireLogin(animeHandler.Edit))
	mux.HandleFunc("PUT /anime/{id}", handlers.RequireLogin(animeHandler.Update))
	mux.HandleFunc("DELETE /anime/{id}", handlers.RequireLogin(animeHandler.Delete))

	// Studio routes
	mux.HandleFunc("GET /studios", handlers.RequireLogin(studioHandler.List))
	mux.HandleFunc("GET /studios/new", handlers.RequireLogin(studioHandler.New))
	mux.HandleFunc("POST /studios", handlers.RequireLogin(studioHandler.Create))
	mux.HandleFunc("DELETE /studios/{id}", handlers.RequireLogin(studioHandler.Delete))

	// User routes (admin only)
	mux.HandleFunc("GET /users", handlers.RequireAdmin(userHandler.List))
	mux.HandleFunc("GET /users/new", handlers.RequireAdmin(userHandler.New))
	mux.HandleFunc("POST /users", handlers.RequireAdmin(userHandler.Create))
	mux.HandleFunc("GET /users/{id}/edit", handlers.RequireAdmin(userHandler.Edit))
	mux.HandleFunc("PUT /users/{id}", handlers.RequireAdmin(userHandler.Update))
	mux.HandleFunc("DELETE /users/{id}", handlers.RequireAdmin(userHandler.Delete))
	mux.HandleFunc("POST /users/{id}/resend-invite", handlers.RequireAdmin(userHandler.ResendInvite))

	// Profile routes
	mux.HandleFunc("GET /profile", handlers.RequireLogin(profileHandler.Page))
	mux.HandleFunc("POST /profile/name", handlers.RequireLogin(profileHandler.UpdateName))
	mux.HandleFunc("POST /profile/password", handlers.RequireLogin(profileHandler.UpdatePassword))

	// Audit routes (admin only)
	mux.HandleFunc("GET /audit-log", handlers.RequireAdmin(auditLogHandler.List))

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// API routes
	mux.HandleFunc("GET /api/anime", handlers.RequireAPIKey(queries)(animeHandler.APIList))

	// Service Account
	mux.HandleFunc("GET /service-accounts", handlers.RequireAdmin(serviceAccountHandler.List))
	mux.HandleFunc("GET /service-accounts/new", handlers.RequireAdmin(serviceAccountHandler.New))
	mux.HandleFunc("POST /service-accounts", handlers.RequireAdmin(serviceAccountHandler.Create))
	mux.HandleFunc("DELETE /service-accounts/{id}", handlers.RequireAdmin(serviceAccountHandler.Delete))

	// Swagger UI
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
