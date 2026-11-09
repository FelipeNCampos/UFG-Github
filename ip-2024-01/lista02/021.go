package main
import ("fmt")

var n,f int 
func main(){
	var ant,pos float64
	for{
		f=0
		fmt.Scan(&n)
		if n==0{
			return
		}
		if n==1{
			fmt.Scan(&ant)
			fmt.Print("Ordenada\n")
			continue
		}
		fmt.Scan(&ant)
		for c:=0;c<n-1;c++{
			fmt.Scan(&pos)
			if ant>pos{
				f = 1
				break
			}
			ant = pos
		}
		if f==0{
			fmt.Print("Ordenada\n")
		}else{
			fmt.Print("Desordenada\n")
		}
	}
}