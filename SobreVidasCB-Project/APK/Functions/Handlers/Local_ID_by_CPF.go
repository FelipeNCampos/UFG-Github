package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"log"
)

func LocalizaIdPorCPF(cpf string) int {
	var id int
	db := Database.OpenConnection()
	row := db.QueryRow("SELECT id FROM ACS WHERE CPF=$1", cpf)
	err := row.Scan(&id)
	if err != nil{
		log.Print("CPF : ",cpf," não cadastrado")
	}
	db.Close()

	return id
}