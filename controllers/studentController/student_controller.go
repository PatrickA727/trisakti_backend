package studentController

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PatrickA727/trisakti-proto/database"
	"github.com/PatrickA727/trisakti-proto/models"
	"github.com/PatrickA727/trisakti-proto/store"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentControllerStruct struct {
	store store.StudentStore
}

func NewController (store store.StudentStore) *StudentControllerStruct {
	return &StudentControllerStruct{
		store: store,
	}
}

func generateFileKey(extension string) string {
	return fmt.Sprintf("students/%s.%s", uuid.New().String(), extension)
}


func (c *StudentControllerStruct) GeneratePresignedUploadURL(ctx *gin.Context) {
	extension := ctx.Query("extension")
	contentType := ctx.Query("content_type")

	fileKey := generateFileKey(extension)

	SignedURL, err := c.store.GeneratePresignedPostURL(fileKey, contentType)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error3": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"presigned_url": SignedURL,
		"fileKey": fileKey,
	})
}

func (c *StudentControllerStruct) GeneratePresignedDownloadURL(ctx *gin.Context) {
	FileKey := ctx.Query("file_key")

	SignedURL, err := c.store.GeneratePresignedGetURL(FileKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error3": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"presigned_url": SignedURL,
	})
}

func (c *StudentControllerStruct) RegisterStudent(ctx *gin.Context) {
	var student models.RegisterStudentPayload;
	var validate = validator.New()

	tx := database.DB.Begin()
	defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

	if err := ctx.ShouldBindJSON(&student); err != nil {  // Binds json payload to the &student model
		ctx.AbortWithStatusJSON(http.StatusBadRequest, 
			gin.H{"error": err.Error()},
		)
		return
	}

	if err := validate.Struct(student); err != nil {
		errors := err.(validator.ValidationErrors)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": errors.Error(),
		})
		return
	}

	if err := c.store.RegisterStudent(tx, &student.NewStudent); err != nil {
		tx.Rollback()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if student.StudentAcademicData != nil {
		for _, acadData := range student.StudentAcademicData {
			acadData.StudentID = student.NewStudent.ID
			if err := c.store.RegisterStudentAcademics(tx, acadData); err != nil {
				tx.Rollback()
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error academic": err.Error(),
				})
			}
		}
	}


	tx.Commit()
	ctx.JSON(http.StatusCreated, "student created")
}

func (c *StudentControllerStruct) GetStudentById(ctx *gin.Context) {
	id := ctx.Param("id");

	id_int, err := strconv.Atoi(id) // Converts string to int
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	student, err := c.store.FindStudentByID(id_int)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H {
			"error": err.Error(),
		})
		return
	}

	acad_data, err := c.store.GetStudentAcademic(id_int)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H {
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"student": student,
		"academic_data": acad_data,
	})
}

func (c *StudentControllerStruct) GetAllStudents(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	satuan := ctx.Query("satuan")
	tahunMasuk := ctx.Query("tahun_masuk")
	nama := ctx.Query("nama")
	limit := 10

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error1": err.Error(),
		})
		return
	}

	offset := (page - 1) * limit

	students, count, err := c.store.GetStudents(offset, limit, satuan, tahunMasuk, nama)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H {
			"error": "Error getting students",
			"message": err.Error(),
		});
		return
	} 

	ctx.JSON(200, gin.H {
		"students": students,
		"count": count,
	})
}

func (c *StudentControllerStruct) UpdateStudentData(ctx *gin.Context) {
	var student models.StudentRegister;
	var updatedStudent models.StudentUpdate;

	// Get student ID from parameter
	id := ctx.Param("id");

	// Check if student exists
	result := database.DB.Table("students").Where("id = ?", id).Find(&student)
	if result.RowsAffected == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Student not found",
		})
		return
	}
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred",
			"error":   result.Error.Error(),
		})
		return
	}

	// Bind new data to updatedStudent
	if err := ctx.ShouldBindJSON(&updatedStudent); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occured",
			"error": err.Error(),
		})
		return
	}

	fmt.Println("DATE: ", updatedStudent.TanggalLahir)

	updatedResult := database.DB.Table("students").Model(&student).Updates(updatedStudent)
	if updatedResult.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occured",
			"error": updatedResult.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"updated student": updatedStudent,
	})

}

func (c *StudentControllerStruct) UpdateStudentAchievment (ctx *gin.Context) {
	var updatedAchievment models.DataAkademikUpdate

	idStr := ctx.Param("id");
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	achievment, err := c.store.GetAchievmentByID(id)
	if err == gorm.ErrRecordNotFound {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Student not found",
		})
		return
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred",
			"error":   err.Error(),
		})
		return
	}

	if err := ctx.ShouldBindJSON(&updatedAchievment); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occured",
			"error": err.Error(),
		})
		return
	}

	err = c.store.UpdateAchievment(achievment, updatedAchievment)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Student achievment updated",
	})
}

func (c *StudentControllerStruct) DeleteStudentAchievment (ctx *gin.Context) {
	idStr := ctx.Param("id");
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	achievment, err := c.store.GetAchievmentByID(id)
	if err == gorm.ErrRecordNotFound {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Student not found",
		})
		return
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred",
			"error":   err.Error(),
		})
		return
	}

	err = c.store.DeleteAchievmentByID(int(achievment.ID))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "data deleted",
	})
}

func (c *StudentControllerStruct) DeleteStudent (ctx *gin.Context) {
	var student models.Students

	id := ctx.Param("id");

	// Delete student
	result := database.DB.Where("id = ?", id).Delete(&student) 

	// Check if student exists
	if result.RowsAffected == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Student not found",
		})
		return
	}
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred",
			"error":   result.Error.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "user deleted",
	})
}

// func (c *StudentControllerStruct) GetSearchedStudents(ctx *gin.Context) {
// 	limitStr := ctx.Query("limit");
// 	offsetStr := ctx.Query("offset");
// 	// searchStr := ctx.Query("search");

// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil || limit <= 0 {
// 		limit = 10
// 	}

// 	offset, err := strconv.Atoi(offsetStr)
// 	if err != nil || offset < 0 {
// 		offset = 0
// 	}
// }

func (c *StudentControllerStruct) GetAllSatuan(ctx *gin.Context) {
	satuan, err := c.store.GetSatuan()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H {
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H {
		"satuan": satuan,
	})
}

func (c *StudentControllerStruct) CreateNewSatuan(ctx *gin.Context) {
	var satuan models.SatuanPendidikanPayload

	if err := ctx.ShouldBindJSON(&satuan); err != nil {  // Binds json payload to the &student model
		ctx.AbortWithStatusJSON(http.StatusBadRequest, 
			gin.H{"error": err.Error()},
		)
		return
	}

	err := c.store.CreateSatuan(models.SatuanPendidikanPayload{
		Satuan: satuan.Satuan,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "created",
	})
}

func (c *StudentControllerStruct) DeleteExistingSatuan(ctx *gin.Context) {
	satuan_id := ctx.Param("id");

	satuan_id_int, err := strconv.Atoi(satuan_id) // Converts string to int
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	err = c.store.DeleteSatuanByID(satuan_id_int)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "deleted",
	})
}

func (c *StudentControllerStruct) CreateNewJurusan(ctx *gin.Context) {
	var jurusan models.Jurusan

	if err := ctx.ShouldBindJSON(&jurusan); err != nil {  // Binds json payload to the &student model
		ctx.AbortWithStatusJSON(http.StatusBadRequest, 
			gin.H{"error": err.Error()},
		)
		return
	}

	err := c.store.CreateJurusan(jurusan)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, 
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(200, gin.H{
		"message": "created",
	})
}

func (c *StudentControllerStruct) DeleteExistingJurusan(ctx *gin.Context) {
	jurusan_id := ctx.Param("id");

	jurusan_id_int, err := strconv.Atoi(jurusan_id) // Converts string to int
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	err = c.store.DeleteJurusanByID(jurusan_id_int)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "deleted",
	})
}

func (c *StudentControllerStruct) GetAllJurusanBySatuan(ctx *gin.Context) {
	satuan_id := ctx.Param("satuan_id");

	satuan_id_int, err := strconv.Atoi(satuan_id) // Converts string to int
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	jurusan, err := c.store.GetJurusanBySatuanID(satuan_id_int)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"jurusan": jurusan,
	})
}

func (c *StudentControllerStruct) CreateStudentAchievment(ctx *gin.Context) {
	var achievment models.DataAkademik

	if err := ctx.ShouldBindJSON(&achievment); err != nil {  // Binds json payload to the &student model
		ctx.AbortWithStatusJSON(http.StatusBadRequest, 
			gin.H{"error": err.Error()},
		)
		return
	}

	err := c.store.CreateAchievment(achievment)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "created",
	})
}

