package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func connect() *sql.DB{
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Erro ao carregar .env!",err)
	}

	usuariobd := os.Getenv("USUARIO")
	psw := os.Getenv("SENHA")
	nomebd := os.Getenv("NOMEBD")
	ddsconnetion := "User="+usuariobd+"dbname="+nomebd+"password="+psw+"host=localhost port=5432 sslmode=disable"
	database,err := sql.Open("postgres",ddsconnetion)
	if err != nil{
		log.Fatal("Erro ao conectar ao banco de dados!",err)
	}
	_,err = database.Query("CREATE TABLE IF NOT EXISTS acess(id SERIAL NOT NULL,pw varchar(30),constraint pk_psw PRIMARY KEY (id));")
	if err != nil{
		log.Fatal("Erro ao iniciar db :",err)
	}
	_,err = database.Query("CREATE TABLE IF NOT EXISTS acs (id SERIAL NOT NULL,nome varchar(50) NOT NULL,cpf bigint UNIQUE NOT NULL,nascimento int NOT NULL,id_gen int NOT NULL,telefone int,email varchar(60),mae varchar(50) NOT NULL,ine int NOT NULL,cbo int NOT NULL,cns bigint NOT NULL,cnes bigint NOT NULL,constraint pk_acs PRIMARY KEY (id),foreign key (id_gen) references genero (id));")
	if err != nil{
		log.Fatal("Erro ao iniciar db :",err)
	}
	_,err = database.Query("CREATE TABLE genero(id SERIAL NOT NULL,nome varchar(30) NOT NULL,CONSTRAINT pk_genero PRIMARY KEY(id));")
	if err != nil{
		log.Fatal("Erro ao iniciar db :",err)
	}
	_,_ = database.Exec("INSERT INTO genero (nome,id) VALUES ('Mulher cis',1);INSERT INTO genero (nome,id) VALUES ('Homem cis',2);INSERT INTO genero (nome,id) VALUES ('Mulher Trans',3);INSERT INTO genero (nome,id) VALUES ('Homem Trans',4);INSERT INTO genero (nome,id) VALUES ('Não-binário',5);INSERT INTO genero (nome,id) VALUES ('Gênero neutro',6);INSERT INTO genero (nome,id) VALUES ('Outros',7);INSERT INTO genero (nome) VALUES ('Não informado');")
	return database
}