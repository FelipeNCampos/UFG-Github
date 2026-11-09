package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"net/http"
)

func SubmitPC(w http.ResponseWriter	, r *http.Request){
	if Database.CheckSesison(){
		if Database.GetSession()==0{
			http.Redirect(w,r,"/",http.StatusTemporaryRedirect)

		}
		if ProcessCadPC(r){
			http.Redirect(w,r,"/resultadopc",http.StatusTemporaryRedirect)
		}else{
			http.Redirect(w,r,"/eresultadopc",http.StatusTemporaryRedirect)
		}

	}else{
		http.Redirect(w,r,"/loginr",http.StatusTemporaryRedirect)
	}
}