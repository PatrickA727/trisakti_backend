package routes

import (
	"github.com/PatrickA727/trisakti-proto/controllers/adminController"
	"github.com/PatrickA727/trisakti-proto/utils"
	"github.com/gin-gonic/gin"
)

func InitAdminRoute(app *gin.Engine, c *adminController.AdminControllerStruct) {
	route := app
	api := route.Group("/api/admin/")
	
	api.POST("/register", c.RegisterAdmin)
	api.POST("/login", c.Login)
	api.POST("/logout", utils.WithJWTAuth(c.Logout, c.Store))
}