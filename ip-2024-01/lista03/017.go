package main
import("fmt")

func main(){
	var x,temp,sum int
	fmt.Scan(&x)
	reg := make([]int,5001)
	for c:=0;c<x;c++{
		fmt.Scan(&temp)
		reg[temp]++
		
	}
	for c:= range reg{
		if reg[c]==1{
			sum++
		}
	}
	fmt.Print(sum)
	
}