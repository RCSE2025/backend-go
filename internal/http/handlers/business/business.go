package business

import (
	"errors"
	"github.com/RCSE2025/backend-go/internal/http/middleware/admin"
	"github.com/RCSE2025/backend-go/internal/http/middleware/logger"
	"github.com/RCSE2025/backend-go/pkg/logger/sl"
	"github.com/gin-contrib/requestid"
	"log/slog"
	"strconv"

	//"github.com/RCSE2025/backend-go/internal/http/middleware/admin"
	"github.com/RCSE2025/backend-go/internal/http/middleware/auth"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/service"
	"github.com/RCSE2025/backend-go/pkg/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type businessRoutes struct {
	s *service.BusinessService
}

func NewBusinessRoutes(h *gin.RouterGroup, s *service.BusinessService, jwtService service.JWTService) {
	g := h.Group("/business")

	ur := businessRoutes{s: s}

	validateJWTmw := auth.ValidateJWT(jwtService)
	onlyAdmin := admin.OnlyAdmin()

	g.GET("/get_business_info/:inn", ur.GetBusinessInfoByINN)

	br := businessRoutes{s: s}

	g.POST("", validateJWTmw, br.CreateBusiness)
	g.GET("/all", validateJWTmw, onlyAdmin, br.GetAllBusinesses)
	g.GET("/:id", validateJWTmw, br.GetBusinessByID)
	g.PUT("/:id", validateJWTmw, br.UpdateBusiness)
	g.DELETE("/:id", validateJWTmw, br.DeleteBusiness)
	g.GET("/inn/:inn", validateJWTmw, br.GetBusinessByINN)
	g.GET("/ogrn/:ogrn", validateJWTmw, br.GetBusinessByOGRN)
	g.GET("/user", validateJWTmw, br.GetUserBusinesses)
	g.GET("/:id/users", validateJWTmw, br.GetBusinessUsers)
	g.POST("/:id/user/:user_id", validateJWTmw, br.AddUserToBusiness)
	g.DELETE("/:id/user/:user_id", validateJWTmw, br.RemoveUserFromBusiness)
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
// @Router      /business/get_business_info/{inn} [get]
// @Param inn path string true "inn"
func (ur *businessRoutes) GetBusinessInfoByINN(c *gin.Context) {
	inn := c.Param("inn")

	resp, err := ur.s.GetBusinessInfoByINN(inn)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid inn"))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateBusiness
// @Summary     Create business
// @Description Create new business
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     400 {object} response.Response
// @Param request body model.Business true "request"
// @Success     201 {object} model.Business
// @Router      /business [post]
// @Security OAuth2PasswordBearer
func (r *businessRoutes) CreateBusiness(c *gin.Context) {
	const op = "handlers.business.CreateBusiness"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	var business model.Business
	if err := c.ShouldBind(&business); err != nil {
		log.Error("cannot parse request", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	userID := c.GetInt64("user_id")
	err := r.s.CreateBusiness(userID, business)
	if err != nil {
		log.Error("cannot create business", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, business)
}

// GetAllBusinesses
// @Summary     Get all businesses
// @Description Get all businesses
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Success     200 {array} model.Business
// @Router      /business/all [get]
// @Security OAuth2PasswordBearer
func (r *businessRoutes) GetAllBusinesses(c *gin.Context) {
	const op = "handlers.business.GetAllBusinesses"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	businesses, err := r.s.GetAllBusinesses()
	if err != nil {
		log.Error("cannot get businesses", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, businesses)
}

// GetBusinessByID
// @Summary     Get business by id
// @Description Get business by id
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param id path string true "id"
// @Success     200 {object} model.Business
// @Router      /business/{id} [get]
// @Security OAuth2PasswordBearer
func (r *businessRoutes) GetBusinessByID(c *gin.Context) {
	const op = "handlers.business.GetBusinessByID"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid id param"))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("cannot parse id"))
		return
	}

	business, err := r.s.GetBusinessByID(idInt)
	if err != nil {
		log.Error("cannot get business", sl.Err(err))
		if errors.Is(err, service.ErrBusinessNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, business)
}

// UpdateBusiness
// @Summary     Update business
// @Description Update business
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param id path string true "id"
// @Param request body model.Business true "request"
// @Success     200 {object} response.Response
// @Router      /business/{id} [put]
// @Security OAuth2PasswordBearer
func (r *businessRoutes) UpdateBusiness(c *gin.Context) {
	const op = "handlers.business.UpdateBusiness"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid id param"))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("cannot parse id"))
		return
	}

	var business model.Business
	if err := c.ShouldBind(&business); err != nil {
		log.Error("cannot parse request", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	err = r.s.UpdateBusiness(idInt, business)
	if err != nil {
		log.Error("cannot update business", sl.Err(err))
		if errors.Is(err, service.ErrBusinessNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("business updated"))
}

// DeleteBusiness
// @Summary     Delete business
// @Description Delete business
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param id path string true "id"
// @Success     200 {object} response.Response
// @Router      /business/{id} [delete]
// @Security OAuth2PasswordBearer
func (r *businessRoutes) DeleteBusiness(c *gin.Context) {
	const op = "handlers.business.DeleteBusiness"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid id param"))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("cannot parse id"))
		return
	}

	err = r.s.DeleteBusiness(idInt)
	if err != nil {
		log.Error("cannot delete business", sl.Err(err))
		if errors.Is(err, service.ErrBusinessNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("business deleted"))
}

// GetBusinessByINN
// @Summary     Get business by INN
// @Description Get business by INN
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param inn path string true "inn"
// @Success     200 {object} model.Business
// @Router      /business/inn/{inn} [get]
// @Security OAuth2PasswordBearer
func (r *businessRoutes) GetBusinessByINN(c *gin.Context) {
	const op = "handlers.business.GetBusinessByINN"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	inn := c.Param("inn")
	if inn == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid inn param"))
		return
	}

	innInt, err := strconv.ParseInt(inn, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("cannot parse inn"))
		return
	}

	business, err := r.s.GetBusinessByINN(innInt)
	if err != nil {
		log.Error("cannot get business", sl.Err(err))
		if errors.Is(err, service.ErrBusinessNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, business)
}

// GetBusinessByOGRN
// @Summary     Get business by OGRN
// @Description Get business by OGRN
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param ogrn path string true "ogrn"
// @Success     200 {object} model.Business
// @Router      /business/ogrn/{ogrn} [get]
// @Security OAuth2PasswordBearer
func (r *businessRoutes) GetBusinessByOGRN(c *gin.Context) {
	const op = "handlers.business.GetBusinessByOGRN"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	ogrn := c.Param("ogrn")
	if ogrn == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid ogrn param"))
		return
	}

	ogrnInt, err := strconv.ParseInt(ogrn, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("cannot parse ogrn"))
		return
	}

	business, err := r.s.GetBusinessByOGRN(ogrnInt)
	if err != nil {
		log.Error("cannot get business", sl.Err(err))
		if errors.Is(err, service.ErrBusinessNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, business)
}

// GetUserBusinesses
// @Summary     Get user businesses
// @Description Get all businesses for a user
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Success     200 {array} model.Business
// @Router      /business/user [get]
// @Security OAuth2PasswordBearer
func (r *businessRoutes) GetUserBusinesses(c *gin.Context) {
	const op = "handlers.business.GetUserBusinesses"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	userID := c.GetInt64("user_id")
	businesses, err := r.s.GetUserBusinesses(userID)
	if err != nil {
		log.Error("cannot get user businesses", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, businesses)
}

// GetBusinessUsers
// @Summary     Get business users
// @Description Get all users for a business
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param id path string true "id"
// @Success     200 {array} model.User
// @Router      /business/{id}/users [get]
// @Security OAuth2PasswordBearer
func (r *businessRoutes) GetBusinessUsers(c *gin.Context) {
	const op = "handlers.business.GetBusinessUsers"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid id param"))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("cannot parse id"))
		return
	}

	users, err := r.s.GetBusinessUsers(idInt)
	if err != nil {
		log.Error("cannot get business users", sl.Err(err))
		if errors.Is(err, service.ErrBusinessNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, users)
}

// AddUserToBusiness
// @Summary     Add user to business
// @Description Add user to business
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param id path string true "business id"
// @Param user_id path string true "user id"
// @Success     200 {object} response.Response
// @Router      /business/{id}/user/{user_id} [post]
// @Security OAuth2PasswordBearer
func (r *businessRoutes) AddUserToBusiness(c *gin.Context) {
	const op = "handlers.business.AddUserToBusiness"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid business id param"))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("cannot parse business id"))
		return
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid user id param"))
		return
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("cannot parse user id"))
		return
	}

	err = r.s.AddUserToBusiness(userIDInt, idInt)
	if err != nil {
		log.Error("cannot add user to business", sl.Err(err))
		if errors.Is(err, service.ErrBusinessNotFound) || errors.Is(err, service.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("user added to business"))
}

// RemoveUserFromBusiness
// @Summary     Remove user from business
// @Description Remove user from business
// @Tags  	    business
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param id path string true "business id"
// @Param user_id path string true "user id"
// @Success     200 {object} response.Response
// @Router      /business/{id}/user/{user_id} [delete]
// @Security OAuth2PasswordBearer
func (r *businessRoutes) RemoveUserFromBusiness(c *gin.Context) {
	const op = "handlers.business.RemoveUserFromBusiness"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid business id param"))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("cannot parse business id"))
		return
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid user id param"))
		return
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("cannot parse user id"))
		return
	}

	err = r.s.RemoveUserFromBusiness(userIDInt, idInt)
	if err != nil {
		log.Error("cannot remove user from business", sl.Err(err))
		if errors.Is(err, service.ErrBusinessNotFound) || errors.Is(err, service.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("user removed from business"))
}
