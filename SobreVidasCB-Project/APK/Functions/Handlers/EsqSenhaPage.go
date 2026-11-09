package Handlers

import (
	//Database "SobreVidasCB-layout/Functions/DB"
	Database "SobreVidasCB-layout/Functions/DB"
	"SobreVidasCB-layout/Functions/Models"
	"net/http"
)


func EsqSenhaPage(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/recuperar" {
		serveTemplate("esqsenha.html",Models.ACS{})(w,r)
		Database.LogOff()
	} else if r.Method == "POST" {
		RecuperarSenha(r)
	}

}