package main
import("fmt")

func fac(n int)int{
	if n==1{
		return n 
	}else{
		return fac(n-1)*n
	}
}

func main(){
	var n int
	fmt.Scan(&n)
	if n<=2{
		fmt.Print("Campeonato inválido !")
		return
	}
	cont:=1
	for c:=0;c<n-1;c++{
		for j:=c+1;j<n;j++{
			fmt.Printf("Final %v:Time%v X Time%v\n",cont,c+1,j+1)
			cont++
		}
	}
}
