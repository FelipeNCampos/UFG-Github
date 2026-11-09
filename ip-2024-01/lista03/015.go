package main
import("fmt"
"math")

func main(){
	var x,temp,min int
	fmt.Scan(&x)
	for c:=0;c<x;c++{
		var tm,cont int
		var reg []int
		fmt.Scan(&tm)
		for j:=0;j<tm;j++{
			fmt.Scan(&temp)
			reg = append(reg,temp)
		}
		min = 99999
		
		for x:=0;x<tm-1;x++{
			for y:=x+1;y<tm;y++{
				if math.Abs(float64(reg[x]-reg[y]))<float64(min){
					min = int(math.Abs(float64(reg[x]-reg[y])))
				}
				cont++
			}
		}
		fmt.Print(min," ",cont,"\n")
	} 
}