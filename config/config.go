package config

import "fmt"

var (
	DBUser     = "clipboard"
	DBPassword = "secret"
	DBName     = "clipboarddb"
	DBHost     = "localhost"
	DBPort     = 5432
)

func GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBName, DBPassword)
}
