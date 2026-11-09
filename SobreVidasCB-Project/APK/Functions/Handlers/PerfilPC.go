package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"net/http"
)

func PerfilPC(w http.ResponseWriter, r *http.Request){
	if Database.CheckSesison(){
		if Database.GetSession()==0{
			http.Redirect(w,r,"/",http.StatusTemporaryRedirect)

		}
		cpf := Database.GetFind()

		data := Get_PC(cpf)
		show := Translate(data)
		serveTemplatepc("perfilpaciente.html",show)(w,r)
	}else{
		http.Redirect(w,r,"/loginr",http.StatusTemporaryRedirect)
	}
}