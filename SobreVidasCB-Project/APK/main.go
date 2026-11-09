package main

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"SobreVidasCB-layout/Functions/Handlers"
	"log"
)

func main(){
	//reinicia a sessão
	Database.SetSession(0)
	Database.LogOff()

	//testa e inicia valores defaults do banco de daods
	db := Database.OpenConnection()
	err := db.Ping()
	if err !=  nil{
		log.Println("Coneção com o banco de dados com problema!\nerr:",err)
	}

	//inicia as rotas
	Handlers.Control()
}