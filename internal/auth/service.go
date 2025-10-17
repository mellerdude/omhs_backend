package auth

import (
	"errors"
	"fmt"
	"time"

	"omhs-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// --- REGISTER ---
func (s *AuthService) Register(req RegisterRequest) (*User, error) {
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, errors.New("all fields are required")
	}

	if existing, _ := s.repo.FindByUsername(req.Username); existing != nil && existing.Username != "" {
		return nil, errors.New("username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:  req.Username,
		Password:  string(hash),
		Email:     req.Email,
		IsAdmin:   false,
		ID:        utils.NewObjectID(),
		LastLogin: time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// --- LOGIN ---
func (s *AuthService) Login(req LoginRequest) (string, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return "", errors.New("invalid username or password")
	}

	token := user.Token
	if time.Since(user.LastLogin) > 48*time.Hour || token == "" {
		token, err = utils.GenerateToken()
		if err != nil {
			return "", err
		}
	}

	user.LastLogin = time.Now()
	if err := s.repo.UpdateToken(user.ID, token, user.LastLogin); err != nil {
		return "", err
	}

	return token, nil
}

// --- RESET PASSWORD ---
func (s *AuthService) ResetPassword(req ResetPasswordRequest) error {
	user, err := s.repo.FindByEmailAndUsername(req.Email, req.Username)
	if err != nil {
		return errors.New("user not found")
	}

	passkey, err := utils.GeneratePasskey()
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePasskey(user.ID, passkey, time.Now()); err != nil {
		return err
	}

	// async invalidation
	go func() {
		time.Sleep(10 * time.Minute)
		_ = s.repo.InvalidatePasskey(user.ID)
	}()

	subject := "Password Reset Passkey"
	message := fmt.Sprintf("Your passkey for resetting your password is: %s", passkey)
	return utils.SendEmail(user.Email, subject, message)
}

// --- CHANGE PASSWORD ---
func (s *AuthService) ChangePassword(req ChangePasswordRequest) error {
	user, err := s.repo.FindByEmailAndUsername(req.Email, req.Username)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Passkey != req.Passkey {
		return errors.New("invalid passkey")
	}
	if time.Since(user.PasskeyGeneratedAt) > 10*time.Minute {
		return errors.New("passkey expired")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(user.ID, string(hash))
}
