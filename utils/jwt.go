package utils

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PatrickA727/trisakti-proto/store"
	// "github.com/PatrickA727/trisakti-proto/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Duration(900) * time.Second	// 15 minutes

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{	// Create new JWt token, with claims(key value pairs embedded in the token)
		"userID": strconv.Itoa(userID),									// Uses the HS256 signing method, its  fast method for single server systems with low complexity
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)	// The final token signed with the secret key
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshJWT(secret []byte, userID int) (string, error) {
	expiration := time.Duration(3600 * 24 * 30) * time.Second // 30 days

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{	// Create new JWt token, with claims(key value pairs embedded in the token)
		"userID": strconv.Itoa(userID),									// Uses the HS256 signing method, its  fast method for single server systems with low complexity
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)	// The final token signed with the secret key
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc gin.HandlerFunc, store store.AdminStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get token from access cookie
		tokenString := GetTokenFromCookie(ctx)

		// Validate JWT
		token, err := ValidateJWT(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error1": err.Error(),
			})
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error2": "invalid token",
			})
			return
		}

		// Get userID from JWT claims
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error3": err.Error(),
			})
			return
		}

		// Fetch user by id from database (Check if user exists)
		u, err := store.GetAdminByID(userID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error4": err.Error(),
			})
			return
		}

		// Set the userId to the ctx(context) so the handler functions have access to current user id in the ctx
		ctx.Set("UserID", u.ID) // Creates a new context that contains UserKey("userid") as the key and user.id as the value

		// Run the handler func with validated user JWT cookie
		handlerFunc(ctx)
	}
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {	// Validates JWT by checking its signing method
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {	// JWT Parse method takes tokenString and a callback func to check/validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {	// Accesses and checks the token signing method (has to be HMAC)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])	// Shows the signing method of the incorrect jwt token
		}

		return []byte(os.Getenv("JWT_SECRET")), nil	// The CALLBACK FUNC returns the secret key to be used by the jwt.parse func
	})
}
