package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"SobreVidasCB-layout/Functions/Models"
	"database/sql"
	"log"
)

func Get_User(id int) Models.ACS{
	db := Database.OpenConnection()
	var data Models.ACS
	rows := db.QueryRow("SELECT nome,cpf,nascimento,id_gen,telefone,email,mae,ine,cbo,cns,cnes from acs where id=$1",id)
	if rows.Err() == sql.ErrNoRows {
		log.Fatal("Sem cadastro\n",rows.Err())
		return Models.ACS{}
	}
	err := rows.Scan(&data.Nome,&data.CPF,&data.Nascimento,&data.Id_Gen,&data.Telefone,&data.Email,&data.Mae,&data.INE,&data.CBO,&data.CNS,&data.CNES)
	if err != nil {
		log.Fatal("Sem cadastro\n",err)
		return Models.ACS{}
	}
	db.Close()
	return data
}