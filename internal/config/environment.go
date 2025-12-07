package config

import(
	"os"
	"log"

	"github.com/joho/godotenv"
)

var Prefix string
var BaseUrl string

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
}