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

func (s *StudentStore) GetStudents(offset int, limit int, satuan string, tahunMasuk string) ([]models.StudentsPayload, error) { 
    var students []models.StudentsPayload
    
    query := s.db.Table("students").Select("students.id", "students.nama", "students.jurusan", "students.tahun_masuk", "students.no_anggota", "students.semester", "satuan_pendidikan.satuan").
                    Joins("LEFT JOIN satuan_pendidikan ON students.satuan_fk = satuan_pendidikan.id").
                    Order("id DESC").Limit(limit).Offset(offset)

	if satuan != "" {
        query = query.Where("satuan_pendidikan.satuan ILIKE ?", satuan+"%")    
    }
    if tahunMasuk != "" {
        query = query.Where("students.tahun_masuk ILIKE ?",  tahunMasuk+"%")
    }

	err := query.Find(&students).Error
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

    result := s.db.Table("students").
        Select("students.*, satuan_pendidikan.satuan").
        Joins("LEFT JOIN satuan_pendidikan ON students.satuan_fk = satuan_pendidikan.id").
        Where("students.id = ?", id).
        First(&student)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("student not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

    return &student, nil
}

func (s *StudentStore) GetSatuan() ([]models.SatuanPendidikan, error) {
    var satuan_list []models.SatuanPendidikan

    err := s.db.Table("satuan_pendidikan").Find(&satuan_list).Error
    if err != nil {
        return nil, err
    }

    return satuan_list, nil
}
