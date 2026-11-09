package main 

import ("fmt")

func main(){
	var temp float64
	fmt.Scan(&temp)
	fmt.Printf("O VALOR EM CELSIUS E = %.2f\n",(5*temp-160)/9)
	fmt.Scan(&temp)
	fmt.Printf("A QUANTIDADE DE CHUVA E = %.2f",temp*25.4)
}
