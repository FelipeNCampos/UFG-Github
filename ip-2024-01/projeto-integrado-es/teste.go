package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgre"
	password = "4321"
	dbname   = "SobreVidasCB"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err!= nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	// Example usage
	addBook(db)
	listBooks(db)
}

func addBook(db *sql.DB) {
	// Placeholder for adding a book
}

func listBooks(db *sql.DB) {
	// Placeholder for listing all books
}

func updateBook(db *sql.DB) {
	// Placeholder for updating a book
}

func deleteBook(db *sql.DB) {
	// Placeholder for deleting a book
}
