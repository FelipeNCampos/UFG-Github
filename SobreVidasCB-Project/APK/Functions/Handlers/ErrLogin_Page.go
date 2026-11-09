package Handlers

import (
	"SobreVidasCB-layout/Functions/Models"
	"net/http"
)

func ErrLogin(w http.ResponseWriter, r *http.Request) {
	serveTemplate("errsenha.html",Models.ACS{})(w,r)
}