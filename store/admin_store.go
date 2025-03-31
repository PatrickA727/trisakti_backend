package store

import (
	"fmt"

	"github.com/PatrickA727/trisakti-proto/models"
	"gorm.io/gorm"
)

type AdminStore struct {
	db *gorm.DB
}

func NewAdminStore (db *gorm.DB) *AdminStore {
	return	&AdminStore{
		db: db,
	}
}

func (s *AdminStore) GetAdminByUname (username string) (*models.AdminUser, error) {
	var admin models.AdminUser

	err := s.db.Table("admin_user").Select("id", "username", "password").Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (s *AdminStore) GetAdminByID (id int) (*models.AdminUser, error) {
	var admin models.AdminUser

	err := s.db.Table("admin_user").Select("id", "username").Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (s *AdminStore) RegisterNewAdmin (admin models.AdminUser) error {
	err := s.db.Create(&admin).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminStore) CreateSession (session models.Sessions) error {
	err := s.db.Create(&session).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminStore) RevokeSession (session models.Sessions) error {
	result := s.db.Model(&session).
        Where("admin_id = ? AND refresh_token = ?", session.AdminID, session.RefreshToken).
        Update("is_revoked", true)

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return fmt.Errorf("session not found or already revoked")
    }

    return nil
}
