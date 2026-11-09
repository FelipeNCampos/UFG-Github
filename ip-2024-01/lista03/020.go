package main
import("fmt"
"math")

func main(){
	var x int
	var temp float64 
	var reg []float64
	fmt.Scan(&x)
	for c:=0;c<3;c++{
		fmt.Scan(&temp)
		reg = append(reg,temp)
	}
	for j :=0;j<x-1;j++{
		var reg2 []float64
		var max float64
		for c:=0;c<3;c++{
			fmt.Scan(&temp)
			reg2 = append(reg2,temp)
			reg[c] = math.Abs(reg[c]-reg2[c])
			if reg[c]>max{
				max = reg[c]
			}
		}
		fmt.Printf("%.2f\n",max)
		for c :=0 ;c<3;c++{
			reg[c] = reg2[c]
		}
	}
}
