package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"SobreVidasCB-layout/Functions/Models"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/login" {
		serveTemplate("login.html",Models.ACS{})(w,r)
		Database.LogOff()
	} else if r.Method == "POST" {
		ProcessLogin(r)
	}

}