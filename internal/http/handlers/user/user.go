package user

import (
	"errors"
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
}

// CreateUser
// @Summary     Create user
// @Description Create new user
// @Tags  	     user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param request body model.UserCreate true "request"
// @Success     201 {object} model.User
// @Router       /user [post]
// @Security Bearer
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
// @Tags  	     user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param id path string true "id"
// @Success     200 {object} model.User
// @Router       /user/{id} [get]
// @Security Bearer
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
// @Tags  	     user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Param id path string true "id"
// @Success     200 {object} response.Response
// @Router       /user/{id} [delete]
// @Security Bearer
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
// @Tags  	     user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Success     200 {array} model.User
// @Router       /user/all [get]
// @Security Bearer
func (r *userRoutes) GetAllUsers(c *gin.Context) {
	users, err := r.s.GetAllUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, users)
}
