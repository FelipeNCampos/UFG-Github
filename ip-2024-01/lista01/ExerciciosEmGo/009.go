package main
import("fmt")

func main(){
	var a,b,c float64
	fmt.Scan(&a)
	fmt.Scan(&b)
	fmt.Scan(&c)
	fmt.Printf("O VALOR DO DELTA E = %.2f\n",((b*b)-(4*a*c)))
}
