package config

import (
	"os"
)

var(
	// Credential DB Main
	DB_MAIN_HOSTNAME	string 
	DB_MAIN_USERNAME	string 
	DB_MAIN_PASSWORD	string 
	DB_MAIN_DBNAME		string 

	// Credential DB MONGO
	DB_MONGO_HOSTNAME	string
	DB_MONGO_USERNAME	string
	DB_MONGO_PASSWORD	string
	DB_MONGO_DBNAME		string
)


func InitDatabaseVars() {
	DB_MAIN_HOSTNAME	= os.Getenv("DB_MAIN_HOSTNAME"+Prefix)
	DB_MAIN_USERNAME	= os.Getenv("DB_MAIN_USERNAME"+Prefix)
	DB_MAIN_PASSWORD	= os.Getenv("DB_MAIN_PASSWORD"+Prefix)
	DB_MAIN_DBNAME		= os.Getenv("DB_MAIN_DBNAME"+Prefix) 

	DB_MONGO_HOSTNAME	= os.Getenv("DB_MONGO_HOSTNAME"+Prefix)
	DB_MONGO_USERNAME	= os.Getenv("DB_MONGO_USERNAME"+Prefix)
	DB_MONGO_PASSWORD	= os.Getenv("DB_MONGO_PASSWORD"+Prefix)
	DB_MONGO_DBNAME		= os.Getenv("DB_MONGO_DBNAME"+Prefix) 
}