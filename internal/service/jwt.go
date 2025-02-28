package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type JWTService interface {
	GenerateToken(userId int64, role string) (string, error)
	GenerateRefreshToken(userId int64, role string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GetUserIDByToken(token string) (int64, error)
	GetUserRole(token string) string
	GenerateRefreshPasswordToken(userId int64) (string, error)
	ValidateRefreshPasswordToken(token string) (refreshPasswordClaim, error)
}

type jwtCustomClaim struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    "Template",
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "Template"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(userId int64, role string) (string, error) {
	claims := jwtCustomClaim{
		userId,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 120)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tx, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return tx, nil
}

func (j *jwtService) parseToken(t_ *jwt.Token) (any, error) {
	if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
	}
	return []byte(j.secretKey), nil
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, j.parseToken)
}

func (j *jwtService) GetUserIDByToken(token string) (int64, error) {
	t_Token, err := j.ValidateToken(token)
	if err != nil {
		return 0, err
	}

	claims := t_Token.Claims.(jwt.MapClaims)
	fmt.Println(claims)
	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found")
	}
	idInt := int64(id)
	return idInt, nil
}

func (j *jwtService) GenerateRefreshToken(userId int64, role string) (string, error) {
	claims := jwtCustomClaim{
		userId,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7 * 30)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tx, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return tx, nil
}

func (j *jwtService) GetUserRole(token string) string {
	t_Token, err := j.ValidateToken(token)
	if err != nil {
		return ""
	}

	claims := t_Token.Claims.(jwt.MapClaims)
	role, ok := claims["role"].(string)
	if !ok {
		return ""
	}
	return role
}

const passwordRefreshTokenExpirationTime = 20 * time.Minute

type refreshPasswordClaim struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (j *jwtService) GenerateRefreshPasswordToken(userId int64) (string, error) {
	claims := refreshPasswordClaim{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(passwordRefreshTokenExpirationTime)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tx, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return tx, nil
}

func (j *jwtService) ValidateRefreshPasswordToken(token string) (refreshPasswordClaim, error) {
	var claims refreshPasswordClaim
	_, err := jwt.ParseWithClaims(token, &claims, j.parseToken)
	if err != nil {
		return refreshPasswordClaim{}, err
	}

	return claims, nil
}
