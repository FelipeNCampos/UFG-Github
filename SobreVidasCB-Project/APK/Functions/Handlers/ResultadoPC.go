package Handlers

import (
	"SobreVidasCB-layout/Functions/Models"
	"net/http"
)

func ResultadoPC(w http.ResponseWriter, r *http.Request){
	var data Models.ACS
	serveTemplate("/sucessinsertpc.html",data)(w,r)
}