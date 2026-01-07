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

	// Soal
	manajemenSoal := r.Group("/manajemen-soal")
	{
		manajemenSoal.Use(middleware.JWTAuthMiddleware(), middleware.LogUserActivity())
		bankSoal := manajemenSoal.Group("/bank")
		{
			bankSoal.GET("", controller.GetBankSoal)
			bankSoal.GET("/:bankId", controller.GetBankSoal)
			bankSoal.DELETE("/:bankId", controller.DeleteBankSoal)

			soalInput := &model.BankSoal{}
			bankSoal.PUT("", middleware.InputValidator(soalInput), controller.InsertBankSoal)
			bankSoal.PATCH("", middleware.InputValidator(soalInput), controller.UpdateBankSoal)
		}

		soalUjian := manajemenSoal.Group("/soal")
		{
			soalUjian.GET("", controller.GetSoalList)
			soalUjian.GET("/:bankId", controller.GetSoalList)
			soalUjian.GET("/:bankId/:soalId", controller.GetSoalList)
			soalUjian.DELETE("/:soalId", controller.DeleteSoal)


			soalUjian.POST("", controller.SaveSoal)
			soalUjian.POST("/import-excel-soal",controller.ReadExcelSoal)
			// soalUjian.PATCH("", controller.SaveSoal)
		}
	}

	// Jadwal Ujian
	jadwalUjian := r.Group("/jadwal-ujian")
	{
		jadwalUjian.Use(middleware.JWTAuthMiddleware(), middleware.LogUserActivity())
		jadwalUjian.GET("", controller.GetJadwalUjian)
		jadwalUjian.GET("/:jadwalId", controller.GetJadwalUjian)
		jadwalUjian.DELETE("/:jadwalId", controller.DeleteJadwalUjian)

		jadwalInput := &model.JadwalUjian{}
		jadwalUjian.PUT("", middleware.InputValidator(jadwalInput), controller.InsertJadwalUjian)
		jadwalUjian.PATCH("", middleware.InputValidator(jadwalInput), controller.UpdateJadwalUjian)

		jadwalUjian.POST("/reset-token", controller.ResetTokenJadwalUjian)
	}

	// Route For Ujian Dari Siswa
	pesertaUjian := r.Group("/peserta-ujian")
	{
		pesertaUjian.Use(middleware.JWTAuthMiddleware(), middleware.LogUserActivity())

		pesertaUjian.GET("/list-jadwal-ujian/:siswaId",controller.GetJadwalUjianSiswaToday)
		pesertaUjian.POST("/proses-token-ujian",controller.ProcessTokenExan)
		pesertaUjian.POST("/save-jawaban-peserta",controller.SaveJawabanPeserta)
	}
}