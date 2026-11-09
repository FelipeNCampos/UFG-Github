package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"net/http"
)

func Pacientes(w http.ResponseWriter, r *http.Request){
	if Database.CheckSesison(){
		if Database.GetSession()==0{
			http.Redirect(w,r,"/",http.StatusTemporaryRedirect)

		}

		if r.URL.Path == "/home/pacientes" {
			data := Get_User(Database.GetSession())
			serveTemplate("pacientes.html",data)(w,r)
		} else {
		if r.Method == "POST" {
			ProcesPC(r)
		}
	}
	}else{
		http.Redirect(w,r,"/loginr",http.StatusTemporaryRedirect)
	}
}