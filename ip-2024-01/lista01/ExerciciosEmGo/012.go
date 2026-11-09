package main 
import("fmt")

func main (){
	var temp int
	fmt.Scan(&temp)
	if temp<3{
		fmt.Printf("O VALO A PAGAR E %i",temp*5)
	}else{
		if temp%3==0{
			fmt.Printf("O VALOR A PAGAR E %.2f",float32(temp)/3*10)
		}else{
			fmt.Printf("O VALOR  A PAGAR E %.2f",(float32((temp)/3)+1)*10)
		}
	}
}
