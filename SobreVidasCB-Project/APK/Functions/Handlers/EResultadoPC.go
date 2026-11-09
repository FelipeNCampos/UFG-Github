package Handlers

import (
	"SobreVidasCB-layout/Functions/Models"
	"net/http"
)

func EResultadoPC(w http.ResponseWriter, r *http.Request){
	var data Models.ACS
	serveTemplate("/errinsertpc.html",data)(w,r)

}