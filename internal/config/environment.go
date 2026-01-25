package config

import(
	"os"
	"log"
	"time"

	"github.com/joho/godotenv"
)

var Prefix string
var BaseUrl string
var SoalPath string
var TimeZone *time.Location

func InitEnvronment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file : %v",err)
	}

	if os.Getenv("ENVIRONMENT") == "development"{
		Prefix = "_DEVELOPMENT"
	}else {
		Prefix = "_PRODUCTION"
	}

	BaseUrl = os.Getenv("BASE_URL"+Prefix)
	SoalPath = os.Getenv("SOAL_PATH")

    TimeZone, err = time.LoadLocation("Asia/Jakarta")
    if err != nil {
        panic("failed to load WIB timezone")
    }
}