package main
import("fmt")



func main(){
	var n,k,temp,cont int
	var reg []int
	fmt.Scan(&n)
	fmt.Scan(&k)
	for c:=0;c<n;c++{
		fmt.Scan(&temp)
		reg = append(reg, temp)
		if temp>0{
			continue
		}else{
			cont++
		}
	}
	if cont>=k{
		fmt.Print("NÃO\n")
		for c:=len(reg)-1;c>-1;c--{
			if reg[c]<=0{
				fmt.Print(c+1,"\n")
			}
		}
	}else{
		fmt.Print("SIM")
	}
}