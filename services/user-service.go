package services

import (
	"carpool-backend/models"
	"carpool-backend/utils"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type UserService interface {
	RegisterUser(user *models.User) error
	LoginUser(email, password string) (*models.User, string, error)
	GetUserByID(id int) (*models.User, error)
	UpdateUser(existingUser *models.User, updates map[string]interface{}) error
	DeleteUser(id int) error
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

func (s *userService) RegisterUser(user *models.User) error {
	// Check if email already exists
	var count int64
	if err := s.db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		return errors.New("failed to check email availability")
	}
	if count > 0 {
		return errors.New("email already in use")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	// Insert the user into the database
	if err := s.db.Create(&user).Error; err != nil {
		return errors.New("failed to register user")
	}

	return nil
}

func (s *userService) LoginUser(email, password string) (*models.User, string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Check the password
	if err := utils.CheckPassword(user.Password, password); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.IsDriver)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return &user, token, nil
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
	// Preserve `created_at`
	updates["created_at"] = existingUser.CreatedAt

	// Ensure `updated_at` is always set
	updates["updated_at"] = time.Now()

	// Perform the update using GORM
	if err := s.db.Model(existingUser).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (s *userService) DeleteUser(id int) error {
	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		return errors.New("failed to delete user")
	}
	return nil
}
