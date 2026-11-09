package main
import("fmt")

func main(){
	var temp,qnt,n,flag int
	fmt.Scan(&qnt)
	var reg []int
	for c:=0 ; c<qnt;c++{
		fmt.Scan(&temp)
		reg = append(reg,temp)
		}	
	fmt.Scan(&n)
	for c:=0;c<n;c++{
	    flag = 0
		fmt.Scan(&temp)
		for j:=0;j<qnt;j++{
			if temp==reg[j]{
				flag =1
				fmt.Print("Achei\n")
				break
			}
		}
	if flag==0{
		fmt.Printf("Nao achei\n")
	}
	}
}
