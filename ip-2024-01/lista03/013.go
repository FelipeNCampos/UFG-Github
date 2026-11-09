package main
import("fmt")

func main(){
	var x,temp,max,id int
	fmt.Scan(&x)
	reg := make([]int,101)
	for c:=0;c<x;c++{
		fmt.Scan(&temp)
		reg[temp]++
	}
	for c:= range reg{
		if reg[c]>max{
			max = reg[c]
			id = c
		}
	}
	fmt.Print(id,"\n",max)
	
}