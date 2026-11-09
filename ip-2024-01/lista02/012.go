package main
import(
	"fmt"
	)

func main(){
	var vi,vn,vm,n,max float64
	var perfect [3]float64
	fmt.Scan(&vi)
	fmt.Scan(&vn)
	fmt.Scan(&vm)
	if vm<vn{
		return 
	}
	for c:=vn;c<=vm;c++{
		if c<=vi{
			n = 120+50*(vi-c)
		}else{
			n = 120-60*(c-vi)
		}
		l := c*n-(200+0.05*n)
		if l>=max{
			max = l
			perfect[0] = c
			perfect[1] = n
			perfect[2] = l
		}
		fmt.Printf("V: %.2f, N: %.0f, L: %.2f\n",c,n,l)
	}
	fmt.Printf("Melhor valor final: %.2f\nLucro: %.2f\nNumero de ingressos: %.0f",perfect[0],perfect[2],perfect[1])
}