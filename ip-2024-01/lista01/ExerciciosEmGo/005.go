package main

import ("fmt")

func main(){
	var conta float32 
	var mm,temp float32
	var tipo string
	fmt.Scan(&conta)
	fmt.Scan(&mm)
	fmt.Scan(&tipo)
	if tipo=="R"{
		temp = 5+(0.05*mm)
	}
	if tipo=="C"{
		if mm<=80{
			temp=500
		}else{
			temp = (500+(0.25*mm))
		}
	} 
	if tipo=="I"{
		if mm<=100{
			temp=800
		}else{
			temp = 800 + (0.04*mm)
		}
	}
	fmt.Print("CONTA = ",conta,"\n")
	fmt.Print("VALOR DA CONTA = ",temp,"\n")
	fmt.Print(conta," ",mm," ",tipo,"\n")
}