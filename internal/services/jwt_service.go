package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"invoiceB2B/internal/config"
	"invoiceB2B/internal/models"
	"math/big"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTService interface {
	// Token generation
	GenerateAccessToken(user *models.User) (string, time.Time, error)
	GenerateRefreshToken(user *models.User) (string, time.Time, error)
	GenerateAccessTokenForStaff(staff *models.Staff) (string, time.Time, error)
	GenerateRefreshTokenForStaff(staff *models.Staff) (string, time.Time, error)
	
	// Token validation and refresh
	ValidateToken(tokenString string, isRefreshToken bool) (jwt.MapClaims, error)
	RefreshAccessToken(refreshToken string) (string, time.Time, error)
	
	// Token revocation
	RevokeToken(tokenID string) error
	RevokeTokenWithExpiration(tokenID string, expiration time.Duration) error
	RevokeAllUserTokens(userID string) error
	RevokeAllTokens() error
	IsTokenRevoked(tokenID string) (bool, error)
	
	// Key management
	GetPrivateKey() (*rsa.PrivateKey, error)
	GetPublicKey() (*rsa.PublicKey, error)
	GetJWKS() (map[string]interface{}, error)
	RotateKeys() error
	ScheduleKeyRotation(interval time.Duration) error
}

type jwtService struct {
	cfg *config.Config
	// Redis client for token revocation and key storage
	redisClient *redis.Client
	// Key prefix for Redis storage
	redisKeyPrefix string
	// Context for Redis operations
	ctx context.Context
}

// Redis key constants
const (
	// Key prefix for revoked tokens
	revokedTokenKeyPrefix = "jwt:revoked:"
	// Key for storing the current private key
	privateKeyRedisKey = "jwt:private_key"
	// Key for storing the key ID
	keyIDRedisKey = "jwt:key_id"
	// Default key ID if none exists
	defaultKeyID = "key1"
)

func NewJWTService(cfg *config.Config, redisClient *redis.Client) JWTService {
	return &jwtService{
		cfg:           cfg,
		redisClient:   redisClient,
		redisKeyPrefix: "jwt:",
		ctx:           context.Background(),
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// generateRSAKeyPair generates a new RSA key pair
func (s *jwtService) generateRSAKeyPair() (*rsa.PrivateKey, error) {
	// Generate a new RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}
	
	// Convert private key to PEM format for storage
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	
	// Store the private key in Redis
	err = s.redisClient.Set(s.ctx, privateKeyRedisKey, string(privateKeyPEM), 0).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store private key in Redis: %w", err)
	}
	
	// Generate a key ID if it doesn't exist
	_, err = s.getKeyID()
	if err != nil {
		// Set a default key ID
		err = s.redisClient.Set(s.ctx, keyIDRedisKey, defaultKeyID, 0).Err()
		if err != nil {
			return nil, fmt.Errorf("failed to store key ID in Redis: %w", err)
		}
	}
	
	return privateKey, nil
}

// GetPrivateKey returns the RSA private key for signing tokens
func (s *jwtService) GetPrivateKey() (*rsa.PrivateKey, error) {
	// Try to get the private key from Redis
	privateKeyPEM, err := s.redisClient.Get(s.ctx, privateKeyRedisKey).Result()
	if err == redis.Nil {
		// Key doesn't exist, generate a new one
		return s.generateRSAKeyPair()
	} else if err != nil {
		return nil, fmt.Errorf("failed to get private key from Redis: %w", err)
	}
	
	// Parse the private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing private key")
	}
	
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	
	return privateKey, nil
}

// GetPublicKey returns the RSA public key for verifying tokens
func (s *jwtService) GetPublicKey() (*rsa.PublicKey, error) {
	privateKey, err := s.GetPrivateKey()
	if err != nil {
		return nil, err
	}
	return &privateKey.PublicKey, nil
}

func (s *jwtService) generateToken(user *models.User, expirationTime time.Duration, isRefreshToken bool) (string, time.Time, error) {
	// Get the private key for signing
	privateKey, err := s.GetPrivateKey()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get private key: %w", err)
	}
	
	// Generate a unique token ID (jti)
	tokenID := uuid.New().String()
	
	// Calculate expiration time (reduced for access tokens)
	var expiration time.Time
	if !isRefreshToken {
		// Reduce access token lifetime to 15 minutes for better security
		expiration = time.Now().Add(15 * time.Minute)
	} else {
		expiration = time.Now().Add(expirationTime)
	}
	
	// Calculate not before time (nbf) - token is valid 1 second after issuance
	notBefore := time.Now()
	
	claims := &Claims{
		UserID: strconv.FormatUint(uint64(user.ID), 10),
		Email:  user.Email,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(notBefore),
			Issuer:    "invoice-financing-app",
			Audience:  []string{"invoice-b2b-api"},
			ID:        tokenID, // JWT ID for revocation tracking
		},
	}
	
	if isRefreshToken {
		claims.RegisteredClaims.Subject = "refresh_token"
	} else {
		claims.RegisteredClaims.Subject = "access_token"
	}

	// Create token with RS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	
	// Sign the token with the private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}
	
	return tokenString, expiration, nil
}

func (s *jwtService) generateTokenForStaff(staff *models.Staff, expirationTime time.Duration, isRefreshToken bool) (string, time.Time, error) {
	// Get the private key for signing
	privateKey, err := s.GetPrivateKey()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get private key: %w", err)
	}
	
	// Generate a unique token ID (jti)
	tokenID := uuid.New().String()
	
	// Calculate expiration time (reduced for access tokens)
	var expiration time.Time
	if !isRefreshToken {
		// Reduce access token lifetime to 15 minutes for better security
		expiration = time.Now().Add(15 * time.Minute)
	} else {
		expiration = time.Now().Add(expirationTime)
	}
	
	// Calculate not before time (nbf) - token is valid 1 second after issuance
	notBefore := time.Now()
	
	claims := &Claims{
		UserID: strconv.FormatUint(uint64(staff.ID), 10),
		Email:  staff.Email,
		Role:   "staff",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(notBefore),
			Issuer:    "invoice-financing-app",
			Audience:  []string{"invoice-b2b-api"},
			ID:        tokenID, // JWT ID for revocation tracking
		},
	}
	
	if isRefreshToken {
		claims.RegisteredClaims.Subject = "refresh_token"
	} else {
		claims.RegisteredClaims.Subject = "access_token"
	}

	// Create token with RS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	
	// Sign the token with the private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}
	
	return tokenString, expiration, nil
}

func (s *jwtService) GenerateAccessTokenForStaff(staff *models.Staff) (string, time.Time, error) {
	return s.generateTokenForStaff(staff, s.cfg.JWTRefreshTokenExpirationDays, false)
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
	// Get the public key for verification
	publicKey, err := s.GetPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}
	
	// Define validation options
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))
	
	// Parse and validate the token
	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		// Check for kid (key ID) in header
		if kid, ok := token.Header["kid"].(string); ok {
			// If key ID matches current key ID, use current public key
			currentKeyID, err := s.getKeyID()
			if err == nil && kid == currentKeyID {
				return publicKey, nil
			}
			
			// Check if it's a previous key during rotation grace period
			previousKeyIDKey := s.redisKeyPrefix + "previous_key_id"
			previousKeyID, err := s.redisClient.Get(s.ctx, previousKeyIDKey).Result()
			if err == nil && kid == previousKeyID {
				// Use the previous key for verification during grace period
				// In a real implementation, we would store and retrieve the previous key
				// For simplicity, we'll just use the current key
				return publicKey, nil
			}
			
			return nil, fmt.Errorf("invalid key ID: %s", kid)
		}
		
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Perform additional claims validation
		
		// 1. Validate token type (access or refresh)
		if isRefreshToken {
			if sub, ok := claims["sub"].(string); !ok || sub != "refresh_token" {
				return nil, fmt.Errorf("invalid token subject for refresh token")
			}
		} else {
			if sub, ok := claims["sub"].(string); !ok || sub != "access_token" {
				return nil, fmt.Errorf("invalid token subject for access token")
			}
		}
		
		// 2. Validate issuer
		if iss, ok := claims["iss"].(string); !ok || iss != "invoice-financing-app" {
			return nil, fmt.Errorf("invalid token issuer")
		}
		
		// 3. Validate audience
		if aud, ok := claims["aud"].([]interface{}); !ok || len(aud) == 0 || aud[0].(string) != "invoice-b2b-api" {
			return nil, fmt.Errorf("invalid token audience")
		}
		
		// 4. Validate expiration time (jwt.Parse already does this, but we can add additional checks)
		if exp, ok := claims["exp"].(float64); !ok {
			return nil, fmt.Errorf("invalid expiration claim")
		} else {
			// Add additional check for very old tokens (e.g., more than 30 days)
			expTime := time.Unix(int64(exp), 0)
			if time.Until(expTime) < -30*24*time.Hour {
				return nil, fmt.Errorf("token is too old")
			}
		}
		
		// 5. Validate issued at time
		if iat, ok := claims["iat"].(float64); !ok {
			return nil, fmt.Errorf("invalid issued-at claim")
		} else {
			// Check for tokens issued in the future (clock skew)
			iatTime := time.Unix(int64(iat), 0)
			if iatTime.After(time.Now().Add(5 * time.Minute)) {
				return nil, fmt.Errorf("token issued in the future (clock skew)")
			}
		}
		
		// 6. Check if token is revoked
		if jti, ok := claims["jti"].(string); ok {
			revoked, err := s.IsTokenRevoked(jti)
			if err != nil {
				return nil, fmt.Errorf("failed to check token revocation status: %w", err)
			}
			if revoked {
				return nil, fmt.Errorf("token has been revoked")
			}
		} else {
			return nil, fmt.Errorf("token missing jti claim")
		}
		
		// 7. Check for user-specific revocation
		if userID, ok := claims["user_id"].(string); ok {
			userRevokeKey := fmt.Sprintf("%suser:%s:revoked_before", s.redisKeyPrefix, userID)
			revokedTimeStr, err := s.redisClient.Get(s.ctx, userRevokeKey).Result()
			if err == nil {
				// User has a revocation timestamp
				revokedTime, err := strconv.ParseInt(revokedTimeStr, 10, 64)
				if err == nil {
					if iat, ok := claims["iat"].(float64); ok {
						tokenIssueTime := int64(iat)
						if tokenIssueTime <= revokedTime {
							return nil, fmt.Errorf("all tokens for this user were revoked")
						}
					}
				}
			}
		}
		
		// 8. Check for global revocation
		globalRevokeKey := s.redisKeyPrefix + "global:revoked_before"
		globalRevokedTimeStr, err := s.redisClient.Get(s.ctx, globalRevokeKey).Result()
		if err == nil {
			// There is a global revocation timestamp
			globalRevokedTime, err := strconv.ParseInt(globalRevokedTimeStr, 10, 64)
			if err == nil {
				if iat, ok := claims["iat"].(float64); ok {
					tokenIssueTime := int64(iat)
					if tokenIssueTime <= globalRevokedTime {
						return nil, fmt.Errorf("all tokens were globally revoked")
					}
				}
			}
		}
		
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RevokeToken adds a token to the revocation list in Redis with default expiration
func (s *jwtService) RevokeToken(tokenID string) error {
	// Default to 24 hours if we can't determine the actual expiration
	return s.RevokeTokenWithExpiration(tokenID, 24*time.Hour)
}

// RevokeTokenWithExpiration adds a token to the revocation list with a specific expiration
func (s *jwtService) RevokeTokenWithExpiration(tokenID string, expiration time.Duration) error {
	key := revokedTokenKeyPrefix + tokenID
	err := s.redisClient.Set(s.ctx, key, "revoked", expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

// RevokeAllUserTokens revokes all tokens for a specific user
func (s *jwtService) RevokeAllUserTokens(userID string) error {
	// Create a user-specific revocation marker with a long expiration (30 days)
	userRevokeKey := fmt.Sprintf("%suser:%s:revoked_before", s.redisKeyPrefix, userID)
	currentTime := time.Now().Unix()
	
	// Store the current timestamp as the revocation time
	err := s.redisClient.Set(s.ctx, userRevokeKey, currentTime, 30*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}
	
	return nil
}

// RevokeAllTokens revokes all tokens by rotating keys and setting a global revocation timestamp
func (s *jwtService) RevokeAllTokens() error {
	// First, rotate the keys to invalidate all existing signatures
	if err := s.RotateKeys(); err != nil {
		return fmt.Errorf("failed to rotate keys during global revocation: %w", err)
	}
	
	// Set a global revocation timestamp with a long expiration (30 days)
	globalRevokeKey := s.redisKeyPrefix + "global:revoked_before"
	currentTime := time.Now().Unix()
	
	err := s.redisClient.Set(s.ctx, globalRevokeKey, currentTime, 30*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set global revocation timestamp: %w", err)
	}
	
	return nil
}

// IsTokenRevoked checks if a token has been revoked
func (s *jwtService) IsTokenRevoked(tokenID string) (bool, error) {
	// First check if the specific token is revoked
	key := revokedTokenKeyPrefix + tokenID
	result, err := s.redisClient.Exists(s.ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token revocation status: %w", err)
	}
	
	if result > 0 {
		return true, nil
	}
	
	// If we have claims, check user-specific and global revocations
	// This would require the token claims, but we don't have them here
	// In a real implementation, we would extract the user ID and issued-at time from the claims
	// and check against the user-specific and global revocation timestamps
	
	return false, nil
}

// GetJWKS returns the JSON Web Key Set for public key distribution
func (s *jwtService) GetJWKS() (map[string]interface{}, error) {
	// Get the current public key
	publicKey, err := s.GetPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}
	
	// Get the current key ID
	keyID, err := s.getKeyID()
	if err != nil {
		return nil, fmt.Errorf("failed to get key ID: %w", err)
	}
	
	// Convert RSA public key to JWK format
	n := base64.URLEncoding.EncodeToString(publicKey.N.Bytes())
	e := base64.URLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes())
	
	// Create JWK
	jwk := map[string]interface{}{
		"kty": "RSA",
		"use": "sig",
		"alg": "RS256",
		"kid": keyID,
		"n":   n,
		"e":   e,
	}
	
	// Create JWKS
	jwks := map[string]interface{}{
		"keys": []map[string]interface{}{jwk},
	}
	
	return jwks, nil
}

// RotateKeys generates a new RSA key pair and updates the key ID
func (s *jwtService) RotateKeys() error {
	// Generate a new RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate new RSA key pair: %w", err)
	}
	
	// Convert private key to PEM format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	
	// Generate a new key ID
	keyID := uuid.New().String()
	
	// Store the previous key ID for a grace period
	previousKeyID, err := s.getKeyID()
	if err == nil && previousKeyID != defaultKeyID {
		// Store the previous key ID with a 24-hour expiration for graceful transition
		previousKeyIDKey := s.redisKeyPrefix + "previous_key_id"
		err = s.redisClient.Set(s.ctx, previousKeyIDKey, previousKeyID, 24*time.Hour).Err()
		if err != nil {
			// Log but don't fail the rotation
			fmt.Printf("Warning: Failed to store previous key ID: %v\n", err)
		}
	}
	
	// Store the new private key in Redis
	err = s.redisClient.Set(s.ctx, privateKeyRedisKey, string(privateKeyPEM), 0).Err()
	if err != nil {
		return fmt.Errorf("failed to store private key in Redis: %w", err)
	}
	
	// Store the new key ID in Redis
	err = s.redisClient.Set(s.ctx, keyIDRedisKey, keyID, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to store key ID in Redis: %w", err)
	}
	
	// Store the rotation timestamp
	rotationTimeKey := s.redisKeyPrefix + "last_rotation_time"
	err = s.redisClient.Set(s.ctx, rotationTimeKey, time.Now().Unix(), 0).Err()
	if err != nil {
		// Log but don't fail the rotation
		fmt.Printf("Warning: Failed to store rotation timestamp: %v\n", err)
	}
	
	return nil
}

// ScheduleKeyRotation sets up automatic key rotation at the specified interval
func (s *jwtService) ScheduleKeyRotation(interval time.Duration) error {
	// Store the rotation interval in Redis
	rotationIntervalKey := s.redisKeyPrefix + "rotation_interval"
	err := s.redisClient.Set(s.ctx, rotationIntervalKey, interval.String(), 0).Err()
	if err != nil {
		return fmt.Errorf("failed to store rotation interval: %w", err)
	}
	
	// Start a background goroutine to handle key rotation
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				// Rotate keys
				if err := s.RotateKeys(); err != nil {
					fmt.Printf("Error rotating keys: %v\n", err)
				} else {
					fmt.Printf("Successfully rotated JWT keys at %s\n", time.Now().Format(time.RFC3339))
				}
			case <-s.ctx.Done():
				// Context canceled, stop the goroutine
				return
			}
		}
	}()
	
	return nil
}

// RefreshAccessToken validates a refresh token and generates a new access token
func (s *jwtService) RefreshAccessToken(refreshToken string) (string, time.Time, error) {
	// Validate the refresh token
	claims, err := s.ValidateToken(refreshToken, true)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid refresh token: %w", err)
	}
	
	// Extract user ID and role from claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return "", time.Time{}, fmt.Errorf("invalid user_id claim in refresh token")
	}
	
	role, ok := claims["role"].(string)
	if !ok {
		return "", time.Time{}, fmt.Errorf("invalid role claim in refresh token")
	}
	
	// Check for token reuse by adding the refresh token to a "used" list
	// This helps prevent refresh token replay attacks
	tokenID, ok := claims["jti"].(string)
	if !ok {
		return "", time.Time{}, fmt.Errorf("invalid jti claim in refresh token")
	}
	
	// Mark the refresh token as used
	usedTokenKey := s.redisKeyPrefix + "used_refresh:" + tokenID
	alreadyUsed, err := s.redisClient.Exists(s.ctx, usedTokenKey).Result()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to check if refresh token was already used: %w", err)
	}
	
	if alreadyUsed > 0 {
		// Token has been used before, potential replay attack
		// Revoke all tokens for this user as a security measure
		userID, _ := strconv.ParseUint(userIDStr, 10, 64)
		if err := s.RevokeAllUserTokens(userIDStr); err != nil {
			fmt.Printf("Warning: Failed to revoke all tokens for user %d after refresh token reuse: %v\n", userID, err)
		}
		return "", time.Time{}, fmt.Errorf("refresh token has been used before (potential replay attack)")
	}
	
	// Mark this refresh token as used with an expiration matching the original token
	// Extract expiration time from claims
	expClaim, ok := claims["exp"].(float64)
	if !ok {
		return "", time.Time{}, fmt.Errorf("invalid exp claim in refresh token")
	}
	
	expTime := time.Unix(int64(expClaim), 0)
	expiresIn := time.Until(expTime)
	if expiresIn < 0 {
		return "", time.Time{}, fmt.Errorf("refresh token has expired")
	}
	
	// Store the used token with the same expiration as the original token
	err = s.redisClient.Set(s.ctx, usedTokenKey, "used", expiresIn).Err()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to mark refresh token as used: %w", err)
	}
	
	// Generate a new access token based on the role
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid user ID format: %w", err)
	}
	
	if role == "staff" {
		// Create a minimal staff model with just the ID
		staff := &models.Staff{}
		staff.Model.ID = uint(userID)
		return s.GenerateAccessTokenForStaff(staff)
	} else {
		// Create a minimal user model with just the ID
		user := &models.User{}
		user.Model.ID = uint(userID)
		return s.GenerateAccessToken(user)
	}
}

// getKeyID retrieves the current key ID from Redis
func (s *jwtService) getKeyID() (string, error) {
	keyID, err := s.redisClient.Get(s.ctx, keyIDRedisKey).Result()
	if err == redis.Nil {
		// Key ID doesn't exist, use default
		return defaultKeyID, nil
	} else if err != nil {
		return "", fmt.Errorf("failed to get key ID from Redis: %w", err)
	}
	
	return keyID, nil
}
