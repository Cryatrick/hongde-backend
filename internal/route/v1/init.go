package v1

import(
	"github.com/gin-gonic/gin"

	"hongde_backend/internal/controller"
	"hongde_backend/internal/middleware"
	// "hongde_backend/internal/model"
)

func InitRoutes(r *gin.RouterGroup) {
	// Public Routes
	r.POST("auth/login", controller.Login)
	// r.GET("/logout/:usrId", controller.Logout)
	r.POST("/refresh-access",middleware.LogUserActivity(), controller.RefreshAccessToken)
}