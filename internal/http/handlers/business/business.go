package business

import (
	//"github.com/RCSE2025/backend-go/internal/http/middleware/admin"
	"github.com/RCSE2025/backend-go/internal/http/middleware/auth"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type businessRoutes struct {
	userService *service.UserService
	//business    *service.Business
}

func NewUserRoutes(h *gin.RouterGroup, s *service.UserService, jwtService service.JWTService) {
	g := h.Group("/business")

	ur := businessRoutes{userService: s}

	validateJWTmw := auth.ValidateJWT(jwtService)
	//onlyAdmin := admin.OnlyAdmin()

	g.POST("", validateJWTmw, ur.CreateBusiness)
	g.GET("/get_business_info/:inn", validateJWTmw, ur.GetBusinessInfoByINN)
	g.GET("/:id", validateJWTmw, ur.GetBusinessByID)
	//g.DELETE("/:id", validateJWTmw, onlyAdmin, ur.DeleteBusinessByID)
}

// CreateBusiness
// @Summary     Create business
// @Description Create new business
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param request body model.Business true "request"
// @Success     201 {object} model.Business
// @Router      /business [post]
func (ur *businessRoutes) CreateBusiness(c *gin.Context) {

	c.JSON(http.StatusOK, model.Business{
		ID:        122,
		INN:       0321,
		OGRN:      nil,
		Owner:     nil,
		ShortName: nil,
		FullName:  nil,
		Address:   nil,
	})
}

// GetBusinessInfoByINN
// @Summary     Get business info by INN from api
// @Description Get business info by INN
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Success     201 {object} model.User`
// @Router      /business/get_business_info/{inn} [post]
// @Param inn path string true "inn"
func (ur *businessRoutes) GetBusinessInfoByINN(c *gin.Context) {

}

// GetBusinessByID
// @Summary     Get business by id
// @Description Get business by id
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Success     201 {object} model.User`
// @Router      /business/{id} [get]
// @Param id path string true "id"
func (ur *businessRoutes) GetBusinessByID(c *gin.Context) {

}
