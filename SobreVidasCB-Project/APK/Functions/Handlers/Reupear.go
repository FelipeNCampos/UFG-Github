package Handlers

import (
	"log"
	"net/http"
)

func RecuperarSenha(r *http.Request) {
	r.ParseForm()

	email := r.Form.Get("email")
	log.Print(email)
	id := LocalizaIdPorEmail(email)
	pass := Acess(id)
	mensage := "Dados de acesso :\nsenha:" + pass
	enviarEmail(email, mensage)
}