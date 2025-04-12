package routes

import (
	"github.com/PatrickA727/trisakti-proto/controllers/adminController"
	"github.com/PatrickA727/trisakti-proto/utils"
	"github.com/gin-gonic/gin"
)

func InitAdminRoute(app *gin.Engine, c *adminController.AdminControllerStruct) {
	route := app
	api := route.Group("/api/admin/")
	
	api.POST("/register", utils.WithJWTAuth(c.RegisterAdmin, c.Store))
	api.POST("/login", c.Login)
	api.POST("/logout", utils.WithJWTAuth(c.Logout, c.Store))
	api.GET("/auth-client", utils.WithJWTAuth(c.AuthCheck, c.Store))
	api.POST("/refresh", c.RenewAccessToken)
}