package main
import("fmt")

func count(reg []int,n int)int{
	var cont int
	for c := range reg{
		if n==reg[c]{
			cont++
		}
	}
	return cont
}

func main(){
	var max, n ,acc, temp int 
	for{
		var reg []int
		fmt.Scan(&n)
		if n==0{
			break
		}
		for c:=0;c<n;c++{
			fmt.Scan(&temp)
			reg = append(reg,temp)
			if temp>max{
				max = temp
			}
		}
		for c:=0;c<=max;c++{
			acc += count(reg,c)
			fmt.Printf("(%v) %v\n",c,acc)
		}
    acc = 0
	}
}
