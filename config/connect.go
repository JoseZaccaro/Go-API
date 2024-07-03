package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var MySqlDatabase *sql.DB

func ConnectToDB() {

	// Use the SQLConnStrHandler to get the connection string
	connStr := SQLConnStrHandler()

	// Open a new database connection
	connection, err := sql.Open("mysql", connStr)
	fmt.Println("After SQL open")

	// Use the HandleError function from your utils package
	HandleError("Failed to connect to database: %v", err)

	// Set the global MySqlDatabase variable
	MySqlDatabase = connection

}

func SQLConnStrHandler() string {
	// Load the .env file
	err := godotenv.Load("./.env")
	HandleError("Error loading .env file: %v", err)

	// Construct the connection string
	result := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Make sure the connection string is correct
	fmt.Printf("Connection string: %s\n", result)
	// Return the connection string
	return result
}

func HandleError(message string, err error) {
	if err != nil {
		log.Printf(message, err)
		os.Exit(1)
	}
}

func CloseConnection() {
	MySqlDatabase.Close()
}
