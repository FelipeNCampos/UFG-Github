package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"net/http"
)

func SubmitPass(w http.ResponseWriter, r* http.Request){
	r.ParseForm()
	email := r.Form.Get("email")
	id := LocalizaIdPorEmail(email)
	pass := Acess(id)
	db := Database.OpenConnection()
	var nome string
	row := db.QueryRow("SELECT nome from acs where id=$1",id)
	row.Scan(&nome)
	corpo := "Acesso do funcionário"+nome+"\n senha:"+pass
	enviarEmail(email, corpo)
	http.Redirect(w,r,"/token",http.StatusPermanentRedirect)
	db.Close()

}