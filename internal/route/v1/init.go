package v1

import(
	"github.com/gin-gonic/gin"

	"hongde_backend/internal/controller"
	"hongde_backend/internal/middleware"
	"hongde_backend/internal/model"
)

func InitRoutes(r *gin.RouterGroup) {
	// Public Routes
	r.POST("/auth/login-admin", controller.LoginAdmin)
	r.POST("/auth/login-siswa", controller.LoginSiswa)
	r.POST("/auth/google-login", controller.LoginGoogle)
	r.POST("/refresh-access",middleware.LogUserActivity(), controller.RefreshAccessToken)

	// Siswa
	siswa := r.Group("/siswa")
	{
		siswa.Use(middleware.JWTAuthMiddleware(), middleware.LogUserActivity())
		siswa.GET("", controller.GetSiswa)
		siswa.GET("/:siswaId", controller.GetSiswa)
		siswa.DELETE("/:siswaId", controller.DeleteSiswa)

		siswaInput := &model.SiswaList{}
		siswa.PUT("", middleware.InputValidator(siswaInput), controller.InsertSiswa)
		siswa.PATCH("", middleware.InputValidator(siswaInput), controller.UpdateSiswa)

		siswa.POST("/reset-password", controller.ResetPasswordSiswa)

		siswa.POST("/upload-excel-siswa",controller.ReadExcelSiswa)
	}
}