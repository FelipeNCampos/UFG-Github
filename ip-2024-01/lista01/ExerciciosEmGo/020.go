package main 
import("fmt")

func main(){
	var h,m,s int 
	fmt.Scan(&h)
	fmt.Scan(&m)
	fmt.Scan(&s)
	fmt.Printf("O TEMPO EM SEGUNDOS E = %v",(s+(m*60)+(h*3600)))
}
