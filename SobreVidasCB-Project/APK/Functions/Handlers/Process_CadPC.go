package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"SobreVidasCB-layout/Functions/Models"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// calcularIdade recebe uma data de nascimento no formato DDMMAAAA e retorna a idade
func calcularIdade(dataNascimento string) (int, error) {
    // Verificar se a string tem o formato correto
    if len(dataNascimento) != 8 {
        return 0, fmt.Errorf("Formato_invalido")
    }

    // Extrair dia, mês e ano do parâmetro de nascimento
    dia, err := strconv.Atoi(dataNascimento[0:2])
    if err != nil {
        return 0, err
    }

    mes, err := strconv.Atoi(dataNascimento[2:4])
    if err != nil {
        return 0, err
    }

    ano, err := strconv.Atoi(dataNascimento[4:8])
    if err != nil {
        return 0, err
    }

    // Obter a data atual
    hoje := time.Now()

    // Criar a data de nascimento
    dataNasc := time.Date(ano, time.Month(mes), dia, 0, 0, 0, 0, time.UTC)

    // Calcular a diferença de anos
    idade := hoje.Year() - dataNasc.Year()

    // Ajustar a idade se o aniversário ainda não tiver ocorrido este ano
    if hoje.Month() < dataNasc.Month() || (hoje.Month() == dataNasc.Month() && hoje.Day() < dataNasc.Day()) {
        idade--
    }

    return idade, nil
}

func ProcessCadPC(r *http.Request)bool{
	//request
	r.ParseForm()
	data := Get_User(Database.GetSession())
	//get input
	var pc Models.Paciente
	var err error
	pc.Nome = r.Form.Get("nome")
	pc.Mae = r.Form.Get("mae")
	pc.CPF = UnMask(r.Form.Get("cpf"))
	pc.Nascimento,_ = strconv.Atoi(UnMask(r.Form.Get("nascimento")))
	pc.Email = r.Form.Get("email")
	pc.Telefone,_ = strconv.Atoi(UnMask(r.Form.Get("telefone")))
	pc.UF = r.Form.Get("uf")
	pc.Municipio = r.Form.Get("municipio")
	pc.CEP,_ = strconv.Atoi(UnMask(r.Form.Get("cep")))
	pc.Bairro = r.Form.Get("bairro")
	pc.Num_casa,_ = strconv.Atoi(UnMask(r.Form.Get("numero")))
	pc.Complemento = r.Form.Get("complemento")
	pc.Email = r.Form.Get("email")

	genero := r.Form.Get("genero")
	switch genero{
	case "Não Informado":
		pc.Id_Gen = 8
	case "Mulher Cis":
		pc.Id_Gen = 1
	case "Homem Cis":
		pc.Id_Gen = 2
	case "Mulher Trans":
		pc.Id_Gen = 3
	case "Homem Trans":
		pc.Id_Gen = 4
	case "Não-Binário":
		pc.Id_Gen = 5
	case "Gênero Neutro":
		pc.Id_Gen = 6
	case "Outros":
		pc.Id_Gen = 7
	}
	pc.Idade,err = calcularIdade(strconv.Itoa(pc.Nascimento))
	if err != nil{
		fmt.Print("Erro ao calcular a idade do pociente ",pc.Nome)
	}
	if pc.Idade>=40{
		pc.M40 = 1
	}else{
		pc.M40 = 0
	}

	if (r.Form.Get("etilista")=="on"){
		pc.Etilista = 1
	}else{
		pc.Etilista = 0
	}
	if (r.Form.Get("tabagista")=="on"){
		pc.Tabagista = 1
	}else{
		pc.Tabagista = 0
	}
	if (r.Form.Get("absteista")=="on"){
		pc.Absteista = 1
	}else{
		pc.Absteista = 0
	}
	if (r.Form.Get("rua")=="on"){
		pc.Situacao = "Rua"
	}else{
		pc.Situacao = "Regular"
	}
	if (r.Form.Get("def")=="on"){
		pc.Deficiente = 1
	}else{
		pc.Deficiente = 0
	}
	//insert
	if InsertPC(pc,data.Id){
		fmt.Print("Paciente ",pc.Nome," inserido banco de dados com sucesso!\n")
		return true
	}else{
		fmt.Print("Paciente ",pc.Nome," nao inserido banco de dados ")
		return false
	}
}