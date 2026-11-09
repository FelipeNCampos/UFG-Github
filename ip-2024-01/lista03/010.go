package main
import("fmt")

func main(){
	var x,temp,max,id int
	fmt.Scan(&x)
	reg := make([]int,10)
	for c:=0;c<x;c++{
		fmt.Scan(&temp)
		reg[temp-1]++
		if temp>max{
			max = temp
			id=c
			}
		}
	fmt.Printf("Nota %v, %v vezes\nNota %v, indice %v\n",temp,reg[temp-1],max,id)
}