package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

const otpLength = 6
const otpChars = "0123456789"

type OTPService interface {
	GenerateAndStoreOTP(ctx context.Context, userID string) (string, error)
	VerifyOTP(ctx context.Context, userID string, otp string) (bool, error)
	DeleteOTP(ctx context.Context, userID string) error
	BlacklistToken(ctx context.Context, tokenStr string, expiry time.Duration) error // Added for Logout
	IsTokenBlacklisted(ctx context.Context, tokenStr string) (bool, error)           // Added for RefreshToken
}

type otpService struct {
	rdb                  *redis.Client
	otpExpiration        time.Duration
	tokenBlacklistPrefix string
}

func NewOTPService(rdb *redis.Client, otpExpirationMinutes time.Duration) OTPService {
	return &otpService{
		rdb:                  rdb,
		otpExpiration:        otpExpirationMinutes,
		tokenBlacklistPrefix: "blacklist:token:",
	}
}

func generateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, buffer)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP bytes: %w", err)
	}

	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%len(otpChars)]
	}
	return string(buffer), nil
}

func (s *otpService) GenerateAndStoreOTP(ctx context.Context, userID string) (string, error) {
	otp, err := generateOTP(otpLength)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("otp:%s", userID)
	err = s.rdb.Set(ctx, key, otp, s.otpExpiration).Err()
	if err != nil {
		log.Printf("Failed to store OTP in Redis for user %s: %v", userID, err)
		return "", fmt.Errorf("could not store OTP: %w", err)
	}
	return otp, nil
}

func (s *otpService) VerifyOTP(ctx context.Context, userID, otpToVerify string) (bool, error) {
	key := fmt.Sprintf("otp:%s", userID)
	storedOTP, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // OTP not found (expired or never set)
	}
	if err != nil {
		log.Printf("Failed to retrieve OTP from Redis for user %s: %v", userID, err)
		return false, fmt.Errorf("could not retrieve OTP: %w", err)
	}

	if storedOTP == otpToVerify {
		return true, nil
	}
	return false, nil
}

func (s *otpService) DeleteOTP(ctx context.Context, userID string) error {
	key := fmt.Sprintf("otp:%s", userID)
	err := s.rdb.Del(ctx, key).Err()
	if err != nil && err != redis.Nil { // Ignore if key doesn't exist
		log.Printf("Failed to delete OTP from Redis for user %s: %v", userID, err)
		return fmt.Errorf("could not delete OTP: %w", err)
	}
	return nil
}

func (s *otpService) BlacklistToken(ctx context.Context, token string, duration time.Duration) error {
	key := s.tokenBlacklistPrefix + token
	err := s.rdb.Set(ctx, key, "blacklisted", duration).Err()
	if err != nil {
		log.Printf("Failed to blacklist token %s: %v", token, err)
		return fmt.Errorf("could not blacklist token: %w", err)
	}
	return nil
}

func (s *otpService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := s.tokenBlacklistPrefix + token
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // Not blacklisted
	}
	if err != nil {
		log.Printf("Error checking token blacklist for %s: %v", token, err)
		return false, fmt.Errorf("could not check token blacklist: %w", err)
	}
	return val == "blacklisted", nil
}
