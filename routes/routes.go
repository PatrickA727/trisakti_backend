package routes

import (
	"github.com/PatrickA727/trisakti-proto/controllers/adminController"
	"github.com/PatrickA727/trisakti-proto/controllers/studentController"
	"github.com/PatrickA727/trisakti-proto/utils"
	"github.com/gin-gonic/gin"
)

func InitRoute(app *gin.Engine, c studentController.StudentControllerStruct, a adminController.AdminControllerStruct) {
	route := app
	api := route.Group("/api/student/")
	
	api.GET("/get-all-student", utils.WithJWTAuth(c.GetAllStudents, a.Store));
	api.GET("/get-student/:id", utils.WithJWTAuth(c.GetStudentById, a.Store));
	api.POST("/register-student", utils.WithJWTAuth(c.RegisterStudent, a.Store));
	api.PATCH("/update-student-data/:id", utils.WithJWTAuth(c.UpdateStudentData, a.Store));
	api.DELETE("/delete-student/:id", utils.WithJWTAuth(c.DeleteStudent, a.Store));
	api.GET("/get-satuan", utils.WithJWTAuth(c.GetAllSatuan, a.Store));
	api.POST("/create-satuan", utils.WithJWTAuth(c.CreateNewSatuan, a.Store));
	api.DELETE("/delete-satuan/:id", utils.WithJWTAuth(c.DeleteExistingSatuan, a.Store));
	api.POST("/create-jurusan", utils.WithJWTAuth(c.CreateNewJurusan, a.Store));
	api.DELETE("/delete-jurusan/:id", utils.WithJWTAuth(c.DeleteExistingJurusan, a.Store));
	api.GET("/get-jurusan-by-satuan/:satuan_id", utils.WithJWTAuth(c.GetAllJurusanBySatuan, a.Store));
	api.POST("/create-achievment", utils.WithJWTAuth(c.CreateStudentAchievment, a.Store));
	api.PATCH("/update-achievment/:id", utils.WithJWTAuth(c.UpdateStudentAchievment, a.Store));
	api.DELETE("/delete-achievment/:id", utils.WithJWTAuth(c.DeleteStudentAchievment, a.Store));
	api.GET("/get-presign", utils.WithJWTAuth(c.GeneratePresignedUploadURL, a.Store));
	api.GET("/get-presign-download", utils.WithJWTAuth(c.GeneratePresignedDownloadURL, a.Store));
}
