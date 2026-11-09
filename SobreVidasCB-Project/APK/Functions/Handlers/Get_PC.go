package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"SobreVidasCB-layout/Functions/Models"
	"log"
)

func Get_PC(cpf string) Models.Paciente{
	db := Database.OpenConnection()
	defer db.Close()
	var data Models.Paciente
	cpf = UnMask(cpf)
	rows := db.QueryRow("SELECT nome,cpf,nascimento,idade,id_gen,telefone,email,mae,uf,bairro,num_casa,cep,complemento,municipio,tabagista,def,etilista,m40,absenteista,acsresponsavel,situacao from pacientes where cpf=$1",cpf)
	err := rows.Scan(&data.Nome,&data.CPF,&data.Nascimento,&data.Idade,&data.Id_Gen,&data.Telefone,&data.Email,&data.Mae,&data.UF,&data.Bairro,&data.Num_casa,&data.CEP,&data.Complemento,&data.Municipio,&data.Tabagista,&data.Deficiente,&data.Etilista,&data.M40,&data.Absteista,&data.ACS_responsavel,&data.Situacao)
	if err != nil {
		log.Print("  Sem cadastro\n",err)
		return Models.Paciente{}
	}
	return data
}