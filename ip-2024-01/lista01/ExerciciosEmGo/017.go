package main 
import("fmt")

func main (){
	var x,y int 
	fmt.Scan(&x)
	fmt.Scan(&y)
	if x%2==0{
		for c:=x;c<x+y*2;c+=2{
			fmt.Printf("%v ",c)
		}
	}else{
		fmt.Printf("O PRIMEIRO NUMERO NAO E PAR.")
	}
}
