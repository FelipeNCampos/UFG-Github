package main
import("fmt")

func main (){
	var n int 
	var temp float64
	fmt.Scan(&n)
	for c:=1;c<=n;c++{
		temp = temp + (1/(float64(c)))
	}
	fmt.Printf("%.6f",temp)
}
