package adminController

import (
	"errors"
	"net/http"

	"github.com/PatrickA727/trisakti-proto/models"
	"github.com/PatrickA727/trisakti-proto/store"
	"github.com/PatrickA727/trisakti-proto/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminControllerStruct struct {
	store store.AdminStore
}

func NewAdminController (store store.AdminStore) *AdminControllerStruct {
	return &AdminControllerStruct{
		store: store,
	}
}

func (c *AdminControllerStruct) Login(ctx *gin.Context) {
	var payload models.AdminUser

	if err := ctx.ShouldBindJSON(&payload); err != nil {  
		ctx.AbortWithStatusJSON(http.StatusBadRequest, 
			gin.H{"error": err.Error()},
		)
		return
	}


}

func (c *AdminControllerStruct) RegisterAdmin(ctx *gin.Context) {
	var payload models.AdminUser

	if err := ctx.ShouldBindJSON(&payload); err != nil {  
		ctx.AbortWithStatusJSON(http.StatusBadRequest, 
			gin.H{"error": err.Error()},
		)
		return
	}

	_, err := c.store.GetAdminByUname(payload.Username)
	if err == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": "username already exists",
		})
		return
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error4": err.Error(),
			})
			return
		}
	}

	hashedPass, err := utils.HashPass(payload.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	err = c.store.RegisterNewAdmin(models.AdminUser{
		Username: payload.Username,
		Password: hashedPass,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H {
		"message": "user created",
	})
}
