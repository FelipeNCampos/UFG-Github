package Handlers

import (
	"SobreVidasCB-layout/Functions/Models"
	"log"
	"net/http"
	"text/template"
)

func serveTemplatepc(filename string,data Models.PacienteC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./Web/Templates/"+filename)
		if err != nil {
			log.Printf("Erro ao analisar o template %s: %v\n", filename, err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Erro ao executar o template %s: %v\n", filename, err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		}
	}
}