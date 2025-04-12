package store

import (
	"fmt"
    "context"
    "os"
    "time"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/PatrickA727/trisakti-proto/models"
	"gorm.io/gorm"
)

type StudentStore struct {
	db *gorm.DB
    BucketName     string
	S3Client       *s3.Client
	PresignClient  *s3.PresignClient
}

func NewStudentStore (db *gorm.DB) *StudentStore {
    cfg := aws.Config{
		Region: "us-east-1", // not used by Biznet Gio, just needs to be consistent
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			os.Getenv("STORE_ACCESS"),
			os.Getenv("STORE_SECRET"),
			"",
		)),
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           "https://nos.jkt-1.neo.id",
					SigningRegion: "us-east-1",
				}, nil
			}),
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

    return &StudentStore{
        db: db,
        BucketName:    "mainbucket",
		S3Client:      s3Client,
		PresignClient: presignClient,
    }
}

func (s *StudentStore) GetStudents(offset int, limit int, satuan string, tahunMasuk string, nama string) ([]models.StudentsPayload, int, error) { 
    var students []models.StudentsPayload
    var count int
    
    query := s.db.Table("students").Select("students.id", "students.nama", "students.jurusan", "students.tahun_masuk", "students.no_anggota", "students.semester", "satuan_pendidikan.satuan, jurusan.nama_jurusan").
                    Joins("LEFT JOIN satuan_pendidikan ON students.satuan_fk = satuan_pendidikan.id").
                    Joins("LEFT JOIN jurusan ON students.jurusan_fk = jurusan.id").
                    Order("id DESC").Limit(limit).Offset(offset)

	if satuan != "" {
        query = query.Where("satuan_pendidikan.satuan ILIKE ?", satuan+"%")    
    }
    if tahunMasuk != "" {
        query = query.Where("students.tahun_masuk ILIKE ?",  tahunMasuk+"%")
    }
    if nama != "" {
        query = query.Where("students.nama ILIKE ?",  nama+"%")
    }

	err := query.Find(&students).Error
    if err != nil {
        return nil, 0, err
    }

    count, err = s.GetStudentsCount(satuan, tahunMasuk, nama)
    if err != nil {
        return nil, 0, err
    }

    return students, count, nil
}

func (s *StudentStore) GetStudentsCount (satuan string, tahunMasuk string, nama string) (int, error) {
    var count int64
    query := s.db.Table("students").Joins("LEFT JOIN satuan_pendidikan ON students.satuan_fk = satuan_pendidikan.id")

    if satuan != "" {
        query = query.Where("satuan_pendidikan.satuan ILIKE ?", satuan+"%")    
    }
    if tahunMasuk != "" {
        query = query.Where("students.tahun_masuk ILIKE ?",  tahunMasuk+"%")
    }
    if nama != "" {
        query = query.Where("students.nama ILIKE ?",  nama+"%")
    }

    err := query.Count(&count).Error
    if err != nil {
        return 0, err
    }

    return int(count), err
}

func (s *StudentStore) RegisterStudent(tx *gorm.DB, student *models.StudentRegister) error {
    if err := tx.Table("students").Create(&student).Error; err != nil {
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

func (s *StudentStore) GetStudentAcademic(student_id int) ([]models.DataAkademik, error) {
    var data_akademik []models.DataAkademik

    result := s.db.Table("data_akademik").
    Select("id, student_id, nama_prestasi, content_link, tipe").
    Where("student_id = ?", student_id).
    Find(&data_akademik) // Executes the query and stores results in dataAkademik

    if result.Error != nil {
        return nil, result.Error // Return an error if the query fails
    }

    return data_akademik, nil
}

func (s *StudentStore) FindStudentByID(id int) (*models.Students, error) {
    var student models.Students

    result := s.db.Table("students").
        Select("students.*, satuan_pendidikan.satuan, jurusan.nama_jurusan").
        Joins("LEFT JOIN satuan_pendidikan ON students.satuan_fk = satuan_pendidikan.id").
        Joins("LEFT JOIN jurusan ON students.jurusan_fk = jurusan.id").
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

func (s *StudentStore) CreateSatuan(satuan models.SatuanPendidikanPayload) error {
    err := s.db.Table("satuan_pendidikan").Create(&satuan).Error
    if err != nil {
        return err
    }

    return nil
}

func (s *StudentStore) DeleteSatuanByID(satuan_id int) error {
    var satuan models.SatuanPendidikan
    err := s.db.Where("id = ?", satuan_id).Delete(&satuan).Error
    if err != nil {
        return err
    }

    return nil
}

func (s *StudentStore) CreateJurusan(jurusan models.Jurusan) error {
    err := s.db.Table("jurusan").Create(&jurusan).Error
    if err != nil {
        return err
    }

    return nil
}

func (s *StudentStore) DeleteJurusanByID(jurusan_id int) error {
    var jurusan models.Jurusan
    err := s.db.Where("id = ?", jurusan_id).Delete(&jurusan).Error
    if err != nil {
        return err
    }

    return nil
}

func (s *StudentStore) GetJurusanBySatuanID(satuan_id int) ([]models.JurusanGet, error) {
    var jurusan []models.JurusanGet

    err := s.db.Table("jurusan").Select("id, nama_jurusan").Where("satuan_id = ?", satuan_id).Find(&jurusan).Error
    if err != nil {
        return nil, err
    }

    return jurusan, nil
}

func (s *StudentStore) CreateAchievment(data_akademik models.DataAkademik) error {
    err := s.db.Create(&data_akademik).Error
    if err != nil {
        return err
    }

    return nil
}

func (s *StudentStore) GeneratePresignedPostURL(fileKey string, contentType string) (string, error) {
	req, err := s.PresignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(fileKey),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (s *StudentStore) GeneratePresignedGetURL(fileKey string) (string, error) {
	req, err := s.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(fileKey),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}
