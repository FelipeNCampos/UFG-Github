package Handlers

import (
	"SobreVidasCB-layout/Functions/Models"
	"strconv"
)


func TranslateA(data Models.ACS)Models.ACSS{
	var resp Models.ACSS


	resp.Nome = data.Nome
	resp.CPF = data.CPF[:3]+"."+data.CPF[3:6]+"."+data.CPF[6:9]+"-"+data.CPF[9:]


	resp.Nascimento = strconv.Itoa(data.Nascimento)[:2]+"/"+strconv.Itoa(data.Nascimento)[2:4]+"/"+strconv.Itoa(data.Nascimento)[4:]


	
	resp.Idade,_ = calcularIdade(strconv.Itoa(data.Nascimento))

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
	resp.INE = data.INE
	resp.CBO = data.CBO
	resp.CNS = data.CNS
	resp.CNES = data.CNES
	resp.Sts = data.Sts

	return resp

}