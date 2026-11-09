package main
import("fmt")

func main(){
	var x,temp,ls int 
	var reg []int
	fmt.Scan(&x)
	for c:=0;c<x;c++{
		fmt.Scan(&temp)
		if temp!=ls{
			reg = append(reg, temp)
		}
		ls = temp
	}
	fmt.Print("\n")
	for c:= range reg{
		fmt.Print(reg[c],"\n")
	}
}