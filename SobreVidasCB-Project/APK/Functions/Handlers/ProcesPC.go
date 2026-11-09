package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"net/http"
)


func ProcesPC(r *http.Request){
	//request
	r.ParseForm()

	//registro dos inputs
	cpf := UnMask(r.Form.Get("cpf"))
	Database.SetFind(cpf)
	
}