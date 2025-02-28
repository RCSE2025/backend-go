package user

import (
	"errors"
	"fmt"
	"github.com/RCSE2025/backend-go/internal/http/middleware/admin"
	"github.com/RCSE2025/backend-go/internal/http/middleware/auth"
	"github.com/RCSE2025/backend-go/internal/http/middleware/logger"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/service"
	"github.com/RCSE2025/backend-go/pkg/api/response"
	"github.com/RCSE2025/backend-go/pkg/logger/sl"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

type userRoutes struct {
	s *service.UserService
}

func NewUserRoutes(h *gin.RouterGroup, s *service.UserService, jwtService service.JWTService) {
	g := h.Group("/user")

	ur := userRoutes{s: s}

	jwtMW := auth.ValidateJWT(jwtService)
	adminMW := admin.OnlyAdmin()
	g.POST("", ur.CreateUser)
	g.GET("/self", jwtMW, ur.Self)
	g.GET("/:id", jwtMW, ur.GetUserByID)
	g.DELETE("/:id", jwtMW, ur.DeleteUserByID)
	g.GET("/all", jwtMW, adminMW, ur.GetAllUsers)
	g.POST("/token", ur.Token)
	g.POST("/refresh", ur.RefreshToken)
	g.POST("/email/verify", jwtMW, ur.VerifyEmail)
	g.POST("/password/reset/email", ur.SendResetPasswordEmail)
	g.POST("/password/reset", ur.RefreshPassword)
	g.GET("/email", jwtMW, adminMW, ur.GetUserByEmail)
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
// @Success     201 {object} model.User`
// @Router      /user [post]
func (r *userRoutes) CreateUser(c *gin.Context) {
	const op = "handlers.user.CreateUser"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	var user model.UserCreate
	if err := c.ShouldBind(&user); err != nil {
		log.Error("cannot parse request", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	userDB, err := r.s.CreateUser(user)
	if err != nil {
		log.Error("cannot create user", sl.Err(err))

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
	const op = "handlers.user.GetUserByID"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	id := c.Param("id")
	if id == "" {
		res := response.Error("invalid id param")
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("cannot parse id"))
		return
	}

	user, err := r.s.GetUserByID(int64(idInt))
	if err != nil {
		log.Error("cannot get user", sl.Err(err))
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
	const op = "handlers.user.DeleteUserByID"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

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
		log.Error("cannot delete user", sl.Err(err))
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
	const op = "handlers.user.GetAllUsers"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	users, err := r.s.GetAllUsers()
	if err != nil {
		log.Error("cannot get users", sl.Err(err))
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
	const op = "handlers.user.Token"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	var req TokenRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	fmt.Println(req.Username)
	fmt.Println(req.Password)

	token, err := r.s.GetToken(req.Username, req.Password)
	if err != nil {
		log.Error("cannot get token", sl.Err(err))
		if errors.Is(err, service.ErrWrongEmailOrPassword) {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, token)

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
// @Router      /user/refresh [post]
// @Param       refresh_token formData string true "Refresh token"
func (r *userRoutes) RefreshToken(c *gin.Context) {
	const op = "handlers.user.RefreshToken"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	var req RefreshTokenRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	token, err := r.s.RefreshToken(req.RefreshToken)
	if err != nil {
		log.Error("cannot refresh token", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, token)
}

// VerifyEmail
// @Summary     Verify email
// @Description Verify email
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Success     200 {object} response.Response
// @Router      /user/email/verify [post]
// @Param       code query string true "Verification code"
// @Security OAuth2PasswordBearer
func (r *userRoutes) VerifyEmail(c *gin.Context) {
	const op = "handlers.user.VerifyEmail"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	code := c.Query("code")
	userID := c.GetInt64("user_id")
	if code == "" {
		res := response.Error("invalid code param")
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	err := r.s.VerifyEmail(userID, code)
	if err != nil {
		log.Error("cannot verify email", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("email verified"))
}

// Self
// @Summary     Get user
// @Description Get user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Success     200 {object} model.User
// @Router      /user/self [get]
// @Security OAuth2PasswordBearer
func (r *userRoutes) Self(c *gin.Context) {
	const op = "handlers.user.Self"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)
	userID := c.GetInt64("user_id")
	user, err := r.s.GetUserByID(userID)
	if err != nil {
		log.Error("cannot get user", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUserByEmail
// @Summary     Get user by email
// @Description Get user by email
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Success     200 {object} model.User
// @Router      /user/email [get]
// @Param email query string true "email"
// @Security OAuth2PasswordBearer
func (r *userRoutes) GetUserByEmail(c *gin.Context) {
	const op = "handlers.user.GetUserByEmail"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	email := c.Query("email")
	if email == "" {
		res := response.Error("invalid email param")
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	user, err := r.s.GetUserByEmail(email)
	if err != nil {
		log.Error("cannot get user", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, user)
}

// SendResetPasswordEmail
// @Summary     Send reset password email
// @Description Send reset password email
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Success     200 {object} response.Response
// @Param email query string true "email"
// @Router      /user/password/reset/email [post]
func (r *userRoutes) SendResetPasswordEmail(c *gin.Context) {
	const op = "handlers.user.RefreshPassword"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	email := c.Query("email")
	if email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid email param"))
		return
	}
	err := r.s.SendResetPasswordEmail(email)
	if err != nil {
		log.Error("cannot send reset password email", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("email sent"))
}

type ResetPasswordRequest struct {
	Token    string `json:"token" form:"token" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// RefreshPassword
// @Summary     Refresh password
// @Description Refresh password
// @Tags  	    user
// @Accept      application/x-www-form-urlencoded
// @Produce     json
// @Failure     500 {object} response.Response
// @Failure     404 {object} response.Response
// @Success     200 {object} response.Response
// @Router      /user/password/reset [post]
// @Param       token formData string true "Refresh token"
// @Param       password formData string true "Password"
func (r *userRoutes) RefreshPassword(c *gin.Context) {
	const op = "handlers.user.RefreshPassword"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	var req ResetPasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error("cannot bind request", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	err := r.s.RefreshPassword(req.Token, req.Password)
	if err != nil {
		log.Error("cannot refresh password", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("password refreshed"))
}
