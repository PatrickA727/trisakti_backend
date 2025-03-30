package store

import (
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

	err := s.db.Table("admin_user").Select("id", "username").Where("username = ?", username).First(&admin).Error
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