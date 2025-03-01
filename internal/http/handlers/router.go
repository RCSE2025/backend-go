package handlers

import (
	"github.com/RCSE2025/backend-go/docs"
	"github.com/RCSE2025/backend-go/internal/http/handlers/business"
	"github.com/RCSE2025/backend-go/internal/http/handlers/cart"
	"github.com/RCSE2025/backend-go/internal/http/handlers/order"
	"github.com/RCSE2025/backend-go/internal/http/handlers/product"
	"github.com/RCSE2025/backend-go/internal/http/handlers/user"
	"github.com/RCSE2025/backend-go/internal/http/middleware"
	mwLogger "github.com/RCSE2025/backend-go/internal/http/middleware/logger"
	mvp "github.com/RCSE2025/backend-go/internal/http/middleware/prometheus"
	"github.com/RCSE2025/backend-go/internal/service"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	"log/slog"
)

// NewRouter -.
// @title           Backend API
// @version         1.0
// @description     API for site
// @securityDefinitions.oauth2.password OAuth2PasswordBearer
// @tokenUrl /user/token
// @scope.read Grants read access
// @scope.write Grants write access

func NewRouter(r *gin.Engine, log *slog.Logger, us *service.UserService, jwtService service.JWTService, productService *service.ProductService, cartService *service.CartService, businessService *service.BusinessService) {

	r.Use(requestid.New()) // Equivalent to middleware.RequestID

	r.Use(mwLogger.New(log)) // Logging middleware

	err := r.SetTrustedProxies(nil) //disabled Trusted Proxies
	if err != nil {
		log.Error(err.Error())
		return
	}

	r.Use(middleware.CORSMiddleware())
	r.Use(mvp.NewGinPrometheusMiddleware("user-api"))
	r.Use(gin.Recovery())

	r.GET("/docs", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/docs/index.html")
	})

	r.GET("/docs/*any", func(context *gin.Context) {
		docs.SwaggerInfo.Host = context.Request.Host
		ginSwagger.CustomWrapHandler(&ginSwagger.Config{URL: "/docs/doc.json"}, swaggerFiles.Handler)(context)
	})

	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	r.Use(middleware.RealIPMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	h := r.Group("")

	user.NewUserRoutes(h, us, jwtService)
	product.NewProductRoutes(h, jwtService, productService)
	cart.NewCartRoutes(h, cartService, jwtService)
	order.NewOrderRoutes(h, orderService, jwtService, cartService)
	business.NewBusinessRoutes(h, businessService, jwtService)
}
