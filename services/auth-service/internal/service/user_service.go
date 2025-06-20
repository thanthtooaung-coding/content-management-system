package service

import (
	"errors"
	"go.uber.org/fx"
	"time"

	"github.com/content-management-system/auth-service/internal/model/types"
	"github.com/content-management-system/auth-service/pkg/db"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var Module = fx.Module("service", fx.Provide(NewUserService))

type UserService struct {
	db     *db.DB
	logger *logrus.Logger
}

func NewUserService(db *db.DB, logger *logrus.Logger) *UserService {
	return &UserService{
		db:     db,
		logger: logger,
	}
}

func (s *UserService) GetAllUsers() ([]types.User, error) {
	var users []types.User
	if err := s.db.Conn.Find(&users).Error; err != nil {
		s.logger.WithError(err).Error("Failed to get all users")
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*types.User, error) {
	var user types.User
	if err := s.db.Conn.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		s.logger.WithError(err).Error("Failed to get user by ID")
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	if err := s.db.Conn.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		s.logger.WithError(err).Error("Failed to get user by email")
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CreateUser(username, email, password string) (*types.User, error) {

	var existingUser types.User
	if err := s.db.Conn.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash password")
		return nil, err
	}

	user := types.User{
		ID:       uuid.New(),
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.db.Conn.Create(&user).Error; err != nil {
		s.logger.WithError(err).Error("Failed to create user")
		return nil, err
	}

	return &user, nil
}

func (s *UserService) ValidatePassword(email, password string) (*types.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *UserService) Register(username, email, password string, roleID uint64) (*types.User, error) {
	_, err := s.GetUserByEmail(email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash password")
		return nil, err
	}

	user := &types.User{
		Username:         username,
		Email:            email,
		Password:         string(hashedPassword),
		RoleID:           roleID,
		RegistrationDate: time.Now(),
	}

	if err := s.db.Conn.Create(user).Error; err != nil {
		s.logger.WithError(err).Error("Failed to create user")
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(email, password string) (*types.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
