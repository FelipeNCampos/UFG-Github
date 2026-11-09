package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"log"
)

func Acess(id int) string {
	var pass string
	db := Database.OpenConnection()
	row := db.QueryRow("SELECT pass from acess where id=$1",id)
	err := row.Scan(&pass)
	if err != nil{
		log.Print("ID",id,"SEM REGISTRO DE SENHA")
	}
	db.Close()
	return pass
}