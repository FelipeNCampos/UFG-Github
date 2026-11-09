package main
import("fmt")

func main(){
	var temp1, temp2 float64 
	fmt.Scan(&temp1)
	fmt.Scan(&temp2)
	pi:=3.14159
	fmt.Printf("O VALOR DE CUSTO E = %.2f", (2*(pi*(temp1*temp1))+(2*pi*temp1*temp2))*100)
}
