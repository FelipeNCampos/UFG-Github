package main

import(
	"fmt"

)

func main(){
	var reg []int
	var temp,n int
	fmt.Print("Digite a quantidade de termos do vetor : ")
	fmt.Scan(&n)
	fmt.Print("Digite os termos do vetor ")
	for c:=0;c<n;c++{
		fmt.Scan(&temp)
		reg = append(reg, temp) 
	}
	for c:=0;c<len(reg)-1;c++{
		for j:=c+1;j<len(reg);j++{
			if reg[c]>reg[j]{
				reg[c],reg[j] = reg[j], reg[c]
			}
		}
	}
	fmt.Print("Registro ordenado: \n", reg)
}