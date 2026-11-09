package main
import(
	"fmt"
	)

func somadiv(n int)int{
	soma :=0
	for j:=1;j<n;j++{
		if n%j==0{
			soma += j
		}
	}
	return soma
}

func main(){
	var amigos int
	cont:=0
	fmt.Scan(&amigos)
	for c:=284;cont<amigos;c++{
		id:=1
		for{
			if id>=c{
				break
			}
			if somadiv(c)==id && somadiv(id)==c{
				fmt.Printf("(%v,%v)\n",id,c)
				cont++
				break
			}
			id++
		}

	}
}