package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/febriandani/material-request-system-backend/internal/domain/master/user"
	"github.com/febriandani/material-request-system-backend/pkg/utils"
)

type Service interface {
	Login(username, password, jwtSecret string, jwtExpiration, jwtExpirationRefresh int64) (*user.User, map[string]string, error)
	ValidateRefreshToken(userID int64, refreshToken string) (bool, string, error)
	Logout(userID int64, refreshToken string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Login(username, password, jwtSecret string, jwtExpiration, jwtExpirationRefresh int64) (*user.User, map[string]string, error) {
	var responseUser user.User

	auth, userID, role, err := s.repo.FindByUsername(username)
	if err != nil {
		return &responseUser, nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(auth.Password),
		[]byte(password),
	) != nil {
		return &responseUser, nil, errors.New("invalid credentials")
	}

	accessToken, _ := utils.GenerateAccessToken(userID, role, jwtSecret, jwtExpiration)
	refreshToken, _ := utils.GenerateRefreshToken(userID, jwtSecret, jwtExpirationRefresh)

	hash := utils.HashToken(refreshToken)
	exp := time.Now().Add(14 * 24 * time.Hour)

	err = s.repo.SaveRefreshToken(userID, hash, exp)
	if err != nil {
		return &responseUser, nil, err
	}

	responseUser.ID = userID
	responseUser.FullName = auth.FullName
	responseUser.Username = auth.Username
	responseUser.Role = role
	responseUser.Email = auth.Email
	responseUser.Phone = auth.Phone

	return &responseUser, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func (s *service) ValidateRefreshToken(
	userID int64,
	refreshToken string,
) (bool, string, error) {

	// 1. Hash refresh token
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	// 2. Validate token existence
	ok, err := s.repo.FindValidRefreshToken(
		context.Background(),
		userID,
		tokenHash,
	)
	if err != nil || !ok {
		return false, "", err
	}

	// 3. Get user role
	role, err := s.repo.GetUserRole(
		context.Background(),
		userID,
	)
	if err != nil {
		return false, "", err
	}

	return true, role, nil
}

func (s *service) Logout(userID int64, refreshToken string) error {
	tokenHashBytes := utils.HashToken(refreshToken)
	return s.repo.RevokeRefreshToken(context.Background(), userID, tokenHashBytes)
}
