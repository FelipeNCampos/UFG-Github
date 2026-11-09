package Handlers

import (
	"SobreVidasCB-layout/Functions/Models"
	"fmt"
	"net/http"
	"strconv"
)

func ProcessCadACS(r *http.Request) {
	r.ParseForm()

	var data Models.ACS

	data.Nome = r.Form.Get("nome")
	data.Mae = r.Form.Get("nomemae")
	data.CPF = UnMask(r.Form.Get("cpf"))
	data.Telefone,_ = strconv.Atoi(UnMask(r.Form.Get("telefone")))
	data.Nascimento,_ = strconv.Atoi(UnMask(r.Form.Get("nascimento")))
	data.Email = r.Form.Get("email")
	Pass := r.Form.Get("senha")
	data.INE,_ = strconv.Atoi(r.Form.Get("ine"))
	data.CBO,_ = strconv.Atoi(r.Form.Get("cbo"))
	data.CNS,_ = strconv.Atoi(r.Form.Get("cns"))
	data.CNES,_ = strconv.Atoi(r.Form.Get("cnes"))
	genero := r.Form.Get("genero")
	switch genero{
	case "Não Informado":
		data.Id_Gen = 8
	case "Mulher Cis":
		data.Id_Gen = 1
	case "Homem Cis":
		data.Id_Gen = 2
	case "Mulher Trans":
		data.Id_Gen = 3
	case "Homem Trans":
		data.Id_Gen = 4
	case "Não-Binário":
		data.Id_Gen = 5
	case "Gênero Neutro":
		data.Id_Gen = 6
	case "Outros":
		data.Id_Gen = 7
	}

	fmt.Print(r.Form.Get("telefone"))
	data.Sts = 1
	InsertACS(data,Pass)

}