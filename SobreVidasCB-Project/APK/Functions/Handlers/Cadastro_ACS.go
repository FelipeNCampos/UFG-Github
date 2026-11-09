package Handlers

import (
	"SobreVidasCB-layout/Functions/Models"
	"net/http"
)

func CadastroACS(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/cadastro" {
		serveTemplate("cadastroacs.html",Models.ACS{})(w,r)
	} else if r.Method == "POST" {
		ProcessCadACS(r)
	}
}