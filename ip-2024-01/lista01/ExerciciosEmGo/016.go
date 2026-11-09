package main
import("fmt")

func main (){
	var temp float64
	fmt.Scan(&temp)
	if temp <=300{
		fmt.Printf("SALARIO COM REAJUSTE : %.2f",temp*1.5)
	}else{
		fmt.Printf("SALARIO COM REAJUSTE : %.2f",temp*1.3)
	}
}
