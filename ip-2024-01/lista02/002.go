package main
import("fmt")

func main(){
	var a,b float32
	var anos int 
	fmt.Scan(&a)
	fmt.Scan(&b)
	for {
		a*=1.03
		b*=1.015
		anos++
		if a>b{
			break
		}
	}
	fmt.Printf("Anos = %v\n",anos)
}
