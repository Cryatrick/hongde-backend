package database

import(
	"log"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"hongde_backend/internal/config"
)

var DbMain *sql.DB

func OpenMain() {
	var err error
	DbMain, err = sql.Open("mysql", config.DB_MAIN_USERNAME+":"+config.DB_MAIN_PASSWORD+"@tcp("+config.DB_MAIN_HOSTNAME+")/"+config.DB_MAIN_DBNAME)

	if err != nil {
		log.Fatalf("Failed to connect to DB Main %v", err)
	}
}