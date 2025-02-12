package routes

import (
	"github.com/PatrickA727/trisakti-proto/controllers/studentController"
	"github.com/gin-gonic/gin"
)

func InitRoute(app *gin.Engine, c studentController.StudentControllerStruct) {
	route := app
	api := route.Group("/api/student/")
	
	api.GET("/get-all-student", c.GetAllStudents);
	api.GET("/get-student/:id", c.GetStudentById);
	api.POST("/register-student", c.RegisterStudent);
	api.PATCH("/update-student-data/:id", c.UpdateStudentData);
	api.DELETE("/delete-student/:id", c.DeleteStudent);
}
