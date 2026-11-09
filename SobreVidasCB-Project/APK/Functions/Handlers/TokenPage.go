package Handlers

import (
	"SobreVidasCB-layout/Functions/Models"
	"net/http"
)

func TokenPage(w http.ResponseWriter, r *http.Request){
	serveTemplate("token.html",Models.ACS{})(w,r)
}