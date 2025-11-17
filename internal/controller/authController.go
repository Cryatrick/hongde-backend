package controller

import(
	"time"
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"hongde_backend/internal/config"
	"hongde_backend/internal/middleware"
	"hongde_backend/internal/service"
)

func Login(c *gin.Context) {
	var loginBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
		RememberMe string `json:"remember_me"`
	}
	if c.Bind(&loginBody) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Failed To read Body",
		})
		return
	}
	user,err := service.LoginAdmin(loginBody.Username, loginBody.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status" : http.StatusUnauthorized,
			"error": "Error Getting User",
		})
		return
	}
	if user == nil {
		user,err = service.LoginSiswa(loginBody.Username, loginBody.Password)
		if err != nil || user == nil {
			c.JSON(http.StatusOK, gin.H{
				"status" : http.StatusUnauthorized,
				"error": "Error Getting User",
			})
			return
		}
	}

	// Generate tokens
	accessToken, _ := middleware.GenerateAccessToken(user.UserId)
	refreshToken, err := middleware.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Something Exploded "+err.Error(),
		})
		return
	}
	status := service.UpsertTokenData(user.UserId, bson.M {
		"user_id":user.UserId,
		"last_ip_address":c.ClientIP(),
		"last_user_agent":c.GetHeader("User-Agent"),
		"access_token":accessToken,
		"refresh_token":refreshToken,
		"refresh_token_expired":time.Now().Add(config.RefreshTokenExpiry),
		"last_login":time.Now(),
		"is_valid_token":"y",
		"is_remember_me":loginBody.RememberMe,
	})
	if status == false {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Failed To Save Token To Mongo DB",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status" :http.StatusOK,
		"userId":user.UserId,
		"userName":user.UserNama,
		"accessToken": accessToken,
		"accessTokenExpiresAt":time.Now().Add(config.AccessTokenExpiry),
		"refreshToken": refreshToken,
	})
	return
}

func RefreshAccessToken (c *gin.Context) {
	// userID := c.GetString("userID")
	refreshToken := c.PostForm("refresh_token")
	userID := c.PostForm("user_id")
	if refreshToken == "" || userID == "" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Please Provide Refresh Token And User ID",
		})
		return
	}
	storedToken := service.GetTokenData(bson.M{"user_id":userID,"refresh_token":refreshToken,"is_valid_token":"y"})
	if storedToken == nil {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusUnauthorized,
			"error": "No Token Found",
		})
		return
	}
	extendRefresh := false
	if time.Now().After(storedToken.RefreshTokenExpiredAt){
		if storedToken.IsRememberMe != "y" {
			c.JSON(http.StatusOK, gin.H{
				"status" :http.StatusUnauthorized,
				"storedToken":storedToken,
				"error": "Refresh Token Expired",
			})
			return
		}
		extendRefresh = true
	}
	if extendRefresh == true {
		func() {
			_ = service.UpsertTokenData(storedToken.UserId, bson.M {"refresh_token_expired": time.Now().Add(config.RefreshTokenExpiry)})
			return 
		}()
	}
	// Generate new access token
	newAccessToken, _ := middleware.GenerateAccessToken(storedToken.UserId)

	// Update access token in the database
	_ = service.UpsertTokenData(storedToken.UserId, bson.M {"access_token": newAccessToken})

	c.JSON(http.StatusOK, gin.H{
		"status" :http.StatusOK,
		"accessToken": newAccessToken,
		"accessTokenExpiresAt":time.Now().Add(config.AccessTokenExpiry),
	})
}