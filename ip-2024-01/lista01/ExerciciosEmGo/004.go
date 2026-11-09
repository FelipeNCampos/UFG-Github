package main

import("fmt")


func main(){
	var temp, pk float32 
	fmt.Scan(&temp)
	pk = (temp*7/10)/100
	fmt.Scan(&temp)
	fmt.Print("Custo por kW : R$",pk,"\n")
	fmt.Print("Custo do consumo : R$",pk*temp,"\n")
	fmt.Print("Custo com desconto : R$",(pk*temp)*9/10,"\n")
}