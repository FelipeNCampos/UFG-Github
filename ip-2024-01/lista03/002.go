package main
import("fmt")

func main(){
	var temp, n,num,res int
	var reg []int
	for {
		fmt.Scan(&temp)
		if temp>0 && temp<=1000{
		    break
		}
	}
	for c:=0;c<temp;c++{
		fmt.Scan(&num)
		reg = append(reg,num)
	}
	
	fmt.Scan(&n)
	for c:=0;c<temp;c++{
		if reg[c]>=n{
			res++
		}
	}
	fmt.Printf("%v\n",res)
}
