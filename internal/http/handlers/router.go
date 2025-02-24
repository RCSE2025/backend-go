package handlers

import (
	"github.com/RCSE2025/backend-go/docs"
	mwLogger "github.com/RCSE2025/backend-go/internal/http/middleware/logger"
	mvp "github.com/RCSE2025/backend-go/internal/http/middleware/prometheus"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	"log/slog"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func RealIPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP() // Extracts real IP
		c.Set("real_ip", clientIP)
		c.Next()
	}
}

// NewRouter -.
// @title           Backend API
// @version         1.0
// @description     API for site
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func NewRouter(r *gin.Engine, log *slog.Logger) {

	r.Use(requestid.New()) // Equivalent to middleware.RequestID

	r.Use(mwLogger.New(log)) // Logging middleware

	err := r.SetTrustedProxies(nil) //disabled Trusted Proxies
	if err != nil {
		log.Error(err.Error())
		return
	}

	r.Use(CORSMiddleware())
	r.Use(mvp.NewGinPrometheusMiddleware("user-api"))
	r.Use(gin.Recovery())

	r.GET("/docs/*any", func(context *gin.Context) {
		docs.SwaggerInfo.Host = context.Request.Host
		ginSwagger.CustomWrapHandler(&ginSwagger.Config{URL: "/docs/doc.json"}, swaggerFiles.Handler)(context)
	})

	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	r.Use(RealIPMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.POST("/user", CreateUser)
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
func CreateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hallo, user",
	})
}
