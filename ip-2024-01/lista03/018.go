package main
import("fmt")

func main(){
	var x,temp int
	fmt.Scan(&x)
	for c:=0 ;c<x;c++{
		var s1,s2 int
		var reg[]int
		for d:=0;d<11;d++{
			
			fmt.Scan(&temp)
			reg = append(reg, temp)
		}
		for x:=1;x<=9;x++{
			s1 += reg[x-1]*x
			s2 += reg[x-1]*(10-x)
		}
		if (s1%11==reg[9]||(s1%11==10 && reg[9]==0)) && (s2%11==reg[10]||(s2%11==10 && reg[10]==0)){
			fmt.Print("CPF Valido\n")
		}else{
			fmt.Print("CPF Invalido\n")
		}
	}
}