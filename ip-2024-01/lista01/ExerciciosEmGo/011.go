package main 
import("fmt")

func main (){
	var temp int 
	fmt.Scan(&temp)
	if temp%3==0 && temp%5==0{
		fmt.Printf("O NUMERO E DIVISIVEL")
	}else{
		fmt.Printf("O NUMERO NAO E DIVISIVEL")
	}
}
