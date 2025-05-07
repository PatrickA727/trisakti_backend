package adminController

import (
	"errors"
	"net/http"
	"os"
	"time"
	"github.com/PatrickA727/trisakti-proto/models"
	"github.com/PatrickA727/trisakti-proto/store"
	"github.com/PatrickA727/trisakti-proto/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminControllerStruct struct {
	Store store.AdminStore
}

func NewAdminController (store store.AdminStore) *AdminControllerStruct {
	return &AdminControllerStruct{
		Store: store,
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

	payload.Username = utils.SanitizeInput(payload.Username)

	a, err := c.Store.GetAdminByUname(payload.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": "incorrect password or email",
		})
		return
	}

	if !utils.ComparePasswords(a.Password, []byte(payload.Password)) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": "incorrect password or email",
		})
		return
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := utils.CreateJWT(secret, a.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	refToken, err := utils.CreateRefreshJWT(secret, a.ID) 
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	err = c.Store.CreateSession(models.Sessions{
		AdminID: a.ID,
		RefreshToken: refToken,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(900 * time.Second),
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user logged in",
	})
}

func (c *AdminControllerStruct) RegisterAdmin(ctx *gin.Context) {
	var payload models.AdminUser

	if err := ctx.ShouldBindJSON(&payload); err != nil {  
		ctx.AbortWithStatusJSON(http.StatusBadRequest, 
			gin.H{"error": err.Error()},
		)
		return
	}

	_, err := c.Store.GetAdminByUname(payload.Username)
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

	err = c.Store.RegisterNewAdmin(models.AdminUser{
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

func(c *AdminControllerStruct) Logout(ctx *gin.Context) {
	token, err := ctx.Cookie("refresh_token")
    if err != nil {
        if errors.Is(err, http.ErrNoCookie) {
            ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "no cookie",
			})
            return 
        }

        ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
        return 
    }

	userID, exists := ctx.Get("UserID")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "no user in context",
		})
		return
	}

	intID, ok := userID.(int)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "incorrect data type for id",
		})
		return
	} 

	err = c.Store.RevokeSession(models.Sessions{
		AdminID: intID,
		RefreshToken: token,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
        Path:     "/",
        HttpOnly: true,	
        Expires:  time.Unix(0, 0), 
        MaxAge:   -1,             
        Secure:   true,  
		SameSite: http.SameSiteLaxMode,
	})
		
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
        Path:     "/",
        HttpOnly: true,	
        Expires:  time.Unix(0, 0), 
        MaxAge:   -1,             
        Secure:   true, 
		SameSite: http.SameSiteLaxMode,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"message": "logged out",
	})
}

func (c *AdminControllerStruct) AuthCheck(ctx *gin.Context) {
	token := utils.GetTokenFromCookie(ctx)

	validToken, err := utils.ValidateJWT(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	if !validToken.Valid {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Authorized",
	})
}

func (c *AdminControllerStruct) RenewAccessToken(ctx *gin.Context) {
	token, err := ctx.Cookie("refresh_token")
    if err != nil {
        if errors.Is(err, http.ErrNoCookie) {
            ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "no cookie",
			})
            return 
        }

        ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error1": err.Error(),
		})
        return 
    }

	sessionExists, adminID, err := c.Store.CheckSession(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error2": err.Error(),
		})
		return
	}

	if !sessionExists {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "session does not exist",
		})
		return
	} else {
		secret := []byte(os.Getenv("JWT_SECRET"))
		token, err := utils.CreateJWT(secret, adminID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error3": err.Error(),
			})
			return
		}

		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:     "access_token",
			Value:    token,
			Expires:  time.Now().Add(600 * time.Second),
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "renewed",
	})
}
