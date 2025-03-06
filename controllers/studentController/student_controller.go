package studentController

import (
	"net/http"
	"strconv"

	"github.com/PatrickA727/trisakti-proto/database"
	"github.com/PatrickA727/trisakti-proto/models"
	"github.com/PatrickA727/trisakti-proto/store"
	"github.com/gin-gonic/gin"
)

type StudentControllerStruct struct {
	store store.StudentStore
}

func NewController (store store.StudentStore) *StudentControllerStruct {
	return &StudentControllerStruct{
		store: store,
	}
}

func (c *StudentControllerStruct) RegisterStudent(ctx *gin.Context) {
	var student models.RegisterStudentPayload;

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

	if err := c.store.RegisterStudent(tx, &student.NewStudent); err != nil {
		tx.Rollback()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, acadData := range student.StudentAcademicData {
		acadData.StudentID = student.NewStudent.ID
		if err := c.store.RegisterStudentAcademics(tx, acadData); err != nil {
			tx.Rollback()
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error academic": err.Error(),
			})
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

	ctx.JSON(200, student)
}

func (c *StudentControllerStruct) GetAllStudents(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	jurusan := ctx.Query("jurusan")
	tahunMasuk := ctx.Query("tahun_masuk")
	limit := 10

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error1": err.Error(),
		})
		return
	}

	offset := (page - 1) * limit

	students, err := c.store.GetStudents(offset, limit, jurusan, tahunMasuk)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H {
			"error": "Error getting students",
			"message": err.Error(),
		});
		return
	} 

	ctx.JSON(200, gin.H {
		"students": students,
	})
}

func (c *StudentControllerStruct) UpdateStudentData(ctx *gin.Context) {
	var student models.Students;
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

	updatedResult := database.DB.Model(&student).Updates(updatedStudent)
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

func (c *StudentControllerStruct) GetSearchedStudents(ctx *gin.Context) {
	limitStr := ctx.Query("limit");
	offsetStr := ctx.Query("offset");
	// searchStr := ctx.Query("search");

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
}
