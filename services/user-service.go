package services

import (
	"carpool-backend/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByIdentifier(identifier string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	UpdateUser(existingUser *models.User, updates map[string]interface{}) error
	UpdateUserByIdentifier(identifier string, updates map[string]interface{}) error
	CountUsersByUsername(username string) (int, error)
	DeleteUser(id int) error
	CountUsersByEmailOrPhone(email, phone string) (int, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	return &userService{db: db}
}

func (s *userService) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *userService) GetUserByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ? OR phone = ? OR username = ?", identifier, identifier, identifier).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) CountUsersByEmailOrPhone(email, phone string) (int, error) {
	var count int64
	if err := s.db.Model(&models.User{}).Where("email = ? OR phone = ?", email, phone).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s *userService) CountUsersByUsername(username string) (int, error) {
	var count int64
	if err := s.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := s.db.Select("id, name, email, phone, is_driver, created_at, updated_at").
		First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *userService) UpdateUser(existingUser *models.User, updates map[string]interface{}) error {
	return s.db.Model(existingUser).Updates(updates).Error
}

func (s *userService) UpdateUserByIdentifier(identifier string, updates map[string]interface{}) error {
	return s.db.Model(&models.User{}).Where("email = ? OR phone = ? OR username = ?", identifier, identifier, identifier).Updates(updates).Error
}

func (s *userService) DeleteUser(id int) error {
	return s.db.Delete(&models.User{}, id).Error
}
