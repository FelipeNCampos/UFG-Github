package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"log"
	"net/http"
)


func ProcessLogin(r *http.Request){
	//request
	r.ParseForm()

	//registro dos inputs
	temp := r.Form.Get("cpf")
	cpf := UnMask(temp)
	senha := r.Form.Get("senha")

	// conecta ao banco de dados 
	if len(cpf)==0{
		return
	}
	//localiza id do cpf fornecido
	if senha != Acess(LocalizaIdPorCPF(cpf)){
		log.Print("Tentativa inválida de login com o CPF : ",cpf)
	}else{
		log.Print("Login realizado pelo CPF : ",cpf)
		Database.SetSession(LocalizaIdPorCPF(cpf))
	}
}