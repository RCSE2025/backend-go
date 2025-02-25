package user

import (
	"errors"
	"fmt"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/service"
	"github.com/RCSE2025/backend-go/pkg/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type userRoutes struct {
	s *service.UserService
}

func NewUserRoutes(h *gin.RouterGroup, s *service.UserService) {
	g := h.Group("/user")

	ur := userRoutes{s: s}

	g.POST("", ur.CreateUser)
	g.GET("/:id", ur.GetUserByID)
	g.DELETE("/:id", ur.DeleteUserByID)
	g.GET("/all", ur.GetAllUsers)
	g.POST("/token", ur.Token)
	g.POST("/refresh", ur.RefreshToken)
}

// CreateUser
// @Summary     Create user
// @Description Create new user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param request body model.UserCreate true "request"
// @Success     201 {object} model.User
// @Router      /user [post]
func (r *userRoutes) CreateUser(c *gin.Context) {
	var user model.UserCreate
	if err := c.ShouldBind(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	userDB, err := r.s.CreateUser(user)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			c.AbortWithStatusJSON(http.StatusConflict, response.Error(err.Error()))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, userDB)
}

// GetUserByID
// @Summary     Get user by id
// @Description Get user by id
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param id path string true "id"
// @Success     200 {object} model.User
// @Router      /user/{id} [get]
// @Security OAuth2PasswordBearer
func (r *userRoutes) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		res := response.Error("invalid id")
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	user, err := r.s.GetUserByID(int64(idInt))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUserByID
// @Summary     Delete user by id
// @Description Delete user by id
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param id path string true "id"
// @Success     200 {object} response.Response
// @Router      /user/{id} [delete]
// @Security OAuth2PasswordBearer
func (r *userRoutes) DeleteUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		res := response.Error("invalid id")
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	err = r.s.DeleteUser(int64(idInt))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK())
}

// GetAllUsers
// @Summary     Get all users
// @Description Get all users
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Success     200 {array} model.User
// @Router      /user/all [get]
// @Security OAuth2PasswordBearer
func (r *userRoutes) GetAllUsers(c *gin.Context) {
	users, err := r.s.GetAllUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, users)
}

type TokenRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// Token
// @Summary     Get token
// @Description Get token
// @Tags  	    user
// @Accept      application/x-www-form-urlencoded
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Success     200 {object} model.Token
// @Router      /user/token [post]
// @Param       username formData string true "Email"
// @Param       password formData string true "Password"
func (r *userRoutes) Token(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	fmt.Println(req.Username)
	fmt.Println(req.Password)
	c.JSON(http.StatusOK, model.Token{
		RefreshToken: "RefreshToken",
		AccessToken:  req.Password,
		TokenType:    "Bearer",
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

// RefreshToken
// @Summary     Refresh token
// @Description Refresh token
// @Tags  	    user
// @Accept      application/x-www-form-urlencoded
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Success     200 {object} model.Token
// @Router       /user/refresh [post]
// @Param       refresh_token formData string true "Refresh token"
func (r *userRoutes) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	fmt.Println(req.RefreshToken)
	c.JSON(http.StatusOK, model.Token{
		RefreshToken: req.RefreshToken,
		AccessToken:  "sfsdfdf",
		TokenType:    "Bearer",
	})
}
