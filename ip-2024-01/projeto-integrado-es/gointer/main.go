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
	user     = "postgres"
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
	inserir(db)
}


func inserir(db *sql.DB) {
	var nome string
	err := db.QueryRow("INSERT INTO acs VALUES ('123-456-789-16','Joaozinho','maria','62 99999-9999','joao@gmail.com',1234567890,123456,'1234-1234-1234-123',1234567) RETURNING nome;").Scan(&nome)
	if err!= nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s foi adicionado no banco de dados \n", nome)


}

