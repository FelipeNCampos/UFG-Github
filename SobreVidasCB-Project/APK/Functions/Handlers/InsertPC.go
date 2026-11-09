package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"SobreVidasCB-layout/Functions/Models"
	"database/sql"
)

func InsertPC(data Models.Paciente,resp int)bool{
	db := Database.OpenConnection()
	var teste int 
	db.Query("INSERT INTO pacientes(nome,cpf,nascimento,idade,id_gen,telefone,email,mae,uf,bairro,num_casa,cep,complemento,tabagista,def,etilista,m40,absenteista,acsresponsavel,situacao,municipio) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)",data.Nome,data.CPF,data.Nascimento,data.Idade,data.Id_Gen,data.Telefone,data.Email,data.Mae,data.UF,data.Bairro,data.Num_casa,data.CEP,data.Complemento,data.Tabagista,data.Deficiente,data.Etilista,data.M40,data.Absteista,data.ACS_responsavel,data.Situacao,data.Municipio)
	row := db.QueryRow("SELECT * FROM pacientes WHERE cpf=$1",data.CPF)
	err := row.Scan(&teste)


	db.Close()
	return err!=sql.ErrNoRows
}