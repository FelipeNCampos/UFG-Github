package main

import (
	"fmt"
	"math"
)

func main(){
	var x int
	var temp float64
	var reg []float64
	fmt.Scan(&x)
	for c:=0;c<x;c++{
		for j:=0;j<3;j++{
			fmt.Scan(&temp)
			reg = append(reg, temp)
		}
	}
	for c:=0;c<x*3-3;c+=3{
		fmt.Printf("%.2f\n",math.Sqrt(float64((reg[c]-reg[c+3])*(reg[c]-reg[c+3])+(reg[c+1]-reg[c+4])*(reg[c+1]-reg[c+4])+(reg[c+2]-reg[c+5])*(reg[c+2]-reg[c+5]))))
	}
}