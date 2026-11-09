package main

import (
	"fmt"
	"math"
)

func fac(n int)int{
	sum := 1
	for c:=2;c<=n;c++{
		sum*=c
	}
	return sum
}

func main(){
	var x,res float64

	var vezes int
	fmt.Scan(&x)
	fmt.Scan(&vezes)
	for c:=0;c<=vezes;c++{
		res += (math.Pow(x,float64(c)))/float64(fac(c))
	}
	fmt.Printf("e^%.2f = %.6f",x,res)
}