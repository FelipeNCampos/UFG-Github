package Handlers

import (
	"fmt"
	"net/http"
)

func Control() {
	//Definição do diretório dos parse servidor
    fs := http.FileServer(http.Dir("./Web/Templates"))
    http.Handle("/", fs)

	//Definição das rotas
	http.HandleFunc("/off",Off)
	http.HandleFunc("/login",LoginPage)
	http.HandleFunc("/loginr",ErrLogin)
	http.HandleFunc("/cadastro",CadastroACS)
	http.HandleFunc("/submitlogin",SubmitLogin)
	http.HandleFunc("/submitcadacs",submitcadacs)
	http.HandleFunc("/recuperar",EsqSenhaPage)
	http.HandleFunc("/token",TokenPage)	
	http.HandleFunc("/home",HomeACS)
	http.HandleFunc("/home/cadastro",CadastroPC)
	http.HandleFunc("/submitpc",SubmitPC)
	http.HandleFunc("/submitpass",SubmitPass)
	http.HandleFunc("/home/submit",SubmitPerfilPC)
	http.HandleFunc("/resultadopc",ResultadoPC)
	http.HandleFunc("/eresultadopc",EResultadoPC)
	http.HandleFunc("/home/dashboard",DashBoard)
	http.HandleFunc("/home/pacientes",Pacientes)
	http.HandleFunc("/home/perfil",PerfilACS)
	http.HandleFunc("/home/paciente",PerfilPC)

    fmt.Println(http.ListenAndServe(":8080", nil))
}