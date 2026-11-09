package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"net/http"
)

func CadastroPC(w http.ResponseWriter, r *http.Request){
	if Database.CheckSesison(){
		if Database.GetSession()==0{
			http.Redirect(w,r,"/",http.StatusTemporaryRedirect)

		}
		data := Get_User(Database.GetSession())
		serveTemplate("cadastro.html",data)(w,r)
	}else{
		http.Redirect(w,r,"/loginr",http.StatusTemporaryRedirect)
	}
}