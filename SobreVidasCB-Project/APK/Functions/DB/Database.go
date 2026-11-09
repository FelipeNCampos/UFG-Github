package Database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func OpenConnection() *sql.DB {
	//Connection info
	User := "postgres"
	Name_DataBase := "SobreVidasCB"
	Pass := "2024"
	ConnectInfo := "user=" + User + " dbname=" + Name_DataBase + " password=" + Pass + " host=localhost port=5432 sslmode=disable"
	
	//coneção ao banco de dados com as informações configutradas
	db, err := sql.Open("postgres",ConnectInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Verifica se a conexão foi estabelecida corretamente
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//inicializando obanco de dados
	db.Query("CREATE TABLE IF NOT EXISTS pacientes(Id SERIAL NOT NULL,Nome varchar(150),CPF varchar(20) UNIQUE,Nascimento int NOT NULL,idade int not null,Id_Gen int NOT NULL,Telefone bigint NOT NULL,Email varchar(150) NOT NULL,Mae varchar(150) NOT NULL,uf varchar(20),bairro varchar(50),num_casa int not null,cep bigint not null,complemento varchar(150),municipio varchar(50) NOT NULL,tabagista int NOT NULL,def int NOT NULL,etilista int NOT NULL,m40 int not null,absenteista int not null,acsresponsavel int,situacao varchar(20),CONSTRAINT pk_pc PRIMARY KEY(id));")
	db.Query("CREATE TABLE IF NOT EXISTS ACS(Id SERIAL NOT NULL,Nome varchar(150),CPF varchar(20) UNIQUE,Nascimento int NOT NULL,Id_Gen int NOT NULL,Telefone bigint NOT NULL,Email varchar(150) NOT NULL,Mae varchar(150) NOT NULL,INE int NOT NULL,CBO int NOT NULL,CNS bigint NOT NULL,CNES bigint NOT NULL,STS int NOT NULL,CONSTRAINT pk_acs PRIMARY KEY(id));")
	db.Query("CREATE TABLE IF NOT EXISTS Acess(Id SERIAL,Pass varchar(100),CONSTRAINT pk_acess PRIMARY KEY(Id));")

	//retorno do ponteiro do banco de dados
	return db
}