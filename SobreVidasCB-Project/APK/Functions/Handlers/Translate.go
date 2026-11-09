package Handlers

import (
	"SobreVidasCB-layout/Functions/Models"
	"strconv"
)

func Translate(data Models.Paciente) Models.PacienteC{
	var resp Models.PacienteC

	resp.Nome = data.Nome
	resp.CPF = data.CPF[:3]+"."+data.CPF[3:6]+"."+data.CPF[6:9]+"-"+data.CPF[9:]


	resp.Nascimento = strconv.Itoa(data.Nascimento)[:2]+"/"+strconv.Itoa(data.Nascimento)[2:4]+"/"+strconv.Itoa(data.Nascimento)[4:]



	resp.Idade = data.Idade

	switch data.Id_Gen {
	case 8:
		resp.Gen = "Não Informado"
	case 1:
		resp.Gen = "Mulher Cis"
	case 2:
		resp.Gen = "Homem Cis"
	case 3:
		resp.Gen = "Mulher Trans"
	case 4:
		resp.Gen = "Homem Trans"
	case 5:
		resp.Gen = "Não-Binário"
	case 6:
		resp.Gen = "Gênero Neutro"
	case 7:
		resp.Gen = "Outros"
	}
	

	resp.Telefone = strconv.Itoa(data.Telefone)[5:]+"-"+strconv.Itoa(data.Telefone)[5:]
	resp.Email = data.Email
	resp.Mae = data.Mae
	resp.UF = data.UF
	resp.Bairro = data.Bairro
	resp.Num_casa = data.Num_casa
	resp.CEP = data.CEP
	resp.Complemento = data.Complemento
	resp.Municipio = data.Municipio
	resp.Tabagista = data.Tabagista
	resp.Deficiente = data.Deficiente
	resp.Etilista = data.Etilista
	resp.M40 = data.M40
	resp.Absteista = data.Absteista
	resp.ACS_responsavel = data.ACS_responsavel
	resp.Situacao = data.Situacao

	return resp
}