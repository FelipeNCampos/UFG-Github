package main
import("fmt")

func main(){
	var vez,fim int 
	fmt.Scan(&vez)
	inicio :=1
	for x :=1;x<=vez;x++{
		fmt.Printf("%v*%v*%v = ",x,x,x)
		fim = inicio+(2*(x-1))
		for c:=inicio;c<=fim;c+=2{
			fmt.Printf("%v",c)
			if c!=fim{
				fmt.Print("+")
			}else{
				fmt.Print("\n")
			}
		}
		inicio = fim+2
	
	}
}