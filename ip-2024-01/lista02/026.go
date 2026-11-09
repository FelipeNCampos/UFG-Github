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
		res += (math.Pow(-1,float64(c))*math.Pow(x,(2*float64(c))+1))/float64(fac(2*c+1))
	}
	fmt.Printf("seno(%.2f) = %.6f",x,res)
}