package services

import (
	"fmt"
	"invoiceB2B/internal/config"
	"invoiceB2B/internal/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateAccessToken(user *models.User) (string, time.Time, error)
	GenerateRefreshToken(user *models.User) (string, time.Time, error)
	GenerateAccessTokenForStaff(staff *models.Staff) (string, time.Time, error)
	GenerateRefreshTokenForStaff(staff *models.Staff) (string, time.Time, error)
	ValidateToken(tokenString string, isRefreshToken bool) (jwt.MapClaims, error)
}

type jwtService struct {
	cfg *config.Config
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{cfg: cfg}
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *jwtService) generateToken(user *models.User, expirationTime time.Duration, isRefreshToken bool) (string, time.Time, error) {
	expiration := time.Now().Add(expirationTime)
	claims := &Claims{
		UserID: strconv.FormatUint(uint64(user.ID), 10),
		Email:  user.Email,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "invoice-financing-app",
		},
	}
	if isRefreshToken {
		claims.RegisteredClaims.Subject = "refresh_token"
	} else {
		claims.RegisteredClaims.Subject = "access_token"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, expiration, nil
}

func (s *jwtService) generateTokenForStaff(staff *models.Staff, expirationTime time.Duration, isRefreshToken bool) (string, time.Time, error) {
	expiration := time.Now().Add(expirationTime)
	claims := &Claims{
		UserID: strconv.FormatUint(uint64(staff.ID), 10),
		Email:  staff.Email,
		Role:   "staff",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "invoice-financing-app",
		},
	}
	if isRefreshToken {
		claims.RegisteredClaims.Subject = "refresh_token"
	} else {
		claims.RegisteredClaims.Subject = "access_token"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, expiration, nil
}

func (s *jwtService) GenerateAccessTokenForStaff(staff *models.Staff) (string, time.Time, error) {
	return s.generateTokenForStaff(staff, s.cfg.JWTAccessTokenExpirationMinutes, false)
}

func (s *jwtService) GenerateRefreshTokenForStaff(staff *models.Staff) (string, time.Time, error) {
	return s.generateTokenForStaff(staff, s.cfg.JWTRefreshTokenExpirationDays, true)
}

func (s *jwtService) GenerateAccessToken(user *models.User) (string, time.Time, error) {
	return s.generateToken(user, s.cfg.JWTAccessTokenExpirationMinutes, false)
}

func (s *jwtService) GenerateRefreshToken(user *models.User) (string, time.Time, error) {
	return s.generateToken(user, s.cfg.JWTRefreshTokenExpirationDays, true)
}

func (s *jwtService) ValidateToken(tokenString string, isRefreshToken bool) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if isRefreshToken {
			if sub, ok := claims["sub"].(string); !ok || sub != "refresh_token" {
				return nil, fmt.Errorf("invalid token subject for refresh token")
			}
		} else {
			if sub, ok := claims["sub"].(string); !ok || sub != "access_token" {
				return nil, fmt.Errorf("invalid token subject for access token")
			}
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
