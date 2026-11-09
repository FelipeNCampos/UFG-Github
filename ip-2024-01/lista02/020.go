package main
import("fmt")

func main(){
	var n,ss int
	fmt.Scan(&n)
	fmt.Print(n," = ")
	for c:=1;c<n;c++{
		if n%c==0{
			ss +=c
			if c !=1{
				fmt.Print(" + ")
			}
			fmt.Printf("%v",c)
		}
	}
	if ss == n{
	
		fmt.Printf(" = %v (Numero perfeito)",ss)
	}else{
		fmt.Printf(" = %v (Numero nao e perfeito)",ss)
	}

}