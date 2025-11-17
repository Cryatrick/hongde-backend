package route

import(
	// "os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"hongde_backend/internal/route/v1"
)
func SetupRouter() *gin.Engine {
	// if os.Getenv("ENVIRONMENT") == "production"{
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1Routes := r.Group("/v1")
	{
		v1.InitRoutes(v1Routes)
	}

	return r
}