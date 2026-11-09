package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"net/http"
)

func SubmitPerfilPC(w http.ResponseWriter, r *http.Request){
	if Database.CheckSesison(){
		if Database.GetSession()==0{
			http.Redirect(w,r,"/",http.StatusTemporaryRedirect)

		}
		ProcesPC(r)
		http.Redirect(w,r,"/home/paciente",http.StatusPermanentRedirect)
	}else{
		http.Redirect(w,r,"/loginr",http.StatusTemporaryRedirect)
	}
}