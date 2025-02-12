package store

import (
	"fmt"

	"github.com/PatrickA727/trisakti-proto/models"
	"gorm.io/gorm"
)

type StudentStore struct {
	db *gorm.DB
}

func NewStudentStore (db *gorm.DB) *StudentStore {
    return &StudentStore{
        db: db,
    }
}

func (s *StudentStore) GetAllStudents() ([]models.Students, error) { 
    var students []models.Students

    err := s.db.Table("students").Find(&students).Error
	if err != nil {
		return nil, err
	} 

    return students, nil
}

func (s *StudentStore) RegisterStudent(tx *gorm.DB, student *models.Students) error {
    if err := tx.Create(&student).Error; err != nil {
		return err
	}

    return nil
}

func (s *StudentStore) RegisterStudentAcademics(tx *gorm.DB, data_akademik models.DataAkademik) error {
	if err := tx.Create(&data_akademik).Error; err != nil {
		return err
	}

	return nil
}

func (s *StudentStore) FindStudentByID(id int) (*models.Students, error) {
    var student models.Students

    result := s.db.Table("students").Where("id = ?", id).Find(&student)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("student not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

    return &student, nil
}
