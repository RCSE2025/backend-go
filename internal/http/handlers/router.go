package handlers

import (
	docs "github.com/RCSE2025/backend-go/docs"
	mvp "github.com/RCSE2025/backend-go/internal/http/middleware/prometheus"
	"github.com/RCSE2025/backend-go/pkg/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	mwLogger "github.com/RCSE2025/backend-go/internal/http/middleware/logger"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"

	"net/http"

	"log/slog"
)

// NewRouter -.
// @title           Backend API
// @version         1.0
// @description     API for site
func NewRouter(r chi.Router, log *slog.Logger) {
	r.Use(middleware.RealIP)
	r.Use(mvp.NewPatternMiddleware("user-api"))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Real-IP"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.RequestID)
	r.Use(mwLogger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Handle("/metrics", promhttp.Handler())

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		// redirect
		http.Redirect(w, r, "/docs/index.html", http.StatusMovedPermanently)
	})

	r.Get("/docs/*", func(w http.ResponseWriter, r *http.Request) {
		baseURL := r.Host
		docs.SwaggerInfo.Host = baseURL
		httpSwagger.Handler(
			httpSwagger.URL("/docs/doc.json"), // The URL pointing to API definition
		).ServeHTTP(w, r)
	})

	r.Post("/user", CreateUser)
}

type CreatUserRequest struct {
	Name        string `json:"name" `
	Patronymic  string `json:"patronymic" `
	Surname     string `json:"surname" `
	Email       string `json:"email" `
	Password    string `json:"password" `
	DateOfBirth string `json:"date_of_birth" `
}

// CreateUser
// @Summary     Create user
// @Description Create new user
// @Tags  	     user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param request body CreatUserRequest true "request"
// @Success     201 {object} model.User
// @Router       /user [post]
// @Security Bearer
func CreateUser(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, response.OK())
}
