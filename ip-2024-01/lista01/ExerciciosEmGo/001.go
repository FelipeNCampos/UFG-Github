package main
import("fmt")

func main(){
	var temp,soma  int 
	
	for i:=0;i<3;i++{
		fmt.Scan(&temp)
		soma +=temp
	}
	if soma/3>=6{
		fmt.Print("Aprovado")
	}else{
		fmt.Print("Reprovado")
	}
}

