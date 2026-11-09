package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"SobreVidasCB-layout/Functions/Models"
	"log"
)

func InsertACS(data Models.ACS,pass string) {
	db := Database.OpenConnection()
	//_,err := db.Query("begin transaction;insert into acs(nome,cpf,nascimento,id_gen,telefone,email,mae,ine,cbo,cns,cnes) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,&10,&11);insert into acess(pass) values($12);commit;",data.Nome, data.CPF,data.Nascimento, data.Id_Gen, data.Telefone, data.Email, data.Mae, data.INE, data.CBO, data.CNS, data.CNES,pass)
	//if err != nil{
	//	log.Println("Erro ao inserir o/a ",data.Nome," no banco de dados\nerr:",err)
	//}
	_, err := db.Query("insert into acs(nome,cpf,nascimento,id_gen,telefone,email,mae,ine,cbo,cns,cnes,sts) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",data.Nome, data.CPF,data.Nascimento, data.Id_Gen, data.Telefone, data.Email, data.Mae, data.INE, data.CBO, data.CNS,data.CNES,data.Sts)
	if err != nil {
		log.Fatal("first query error:",err)
	}

	_, err = db.Query("insert into acess(pass) values($1)",pass)
	if err != nil {
		log.Fatal("second query error",err)
	}
	log.Println(data.Nome," foi inserido na tabela de ACS.")
	db.Close()

}	