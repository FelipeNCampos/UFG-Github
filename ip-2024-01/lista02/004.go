package main
import(
	"fmt"
	)

func main(){
	var n,i,s float32
	var k int  
	fmt.Scan(&n)
	fmt.Scan(&i)
	fmt.Scan(&k)
	fmt.Scan(&s)
	fmt.Printf("Tabuada de soma : \n")
	for c:=0;c<k;c++{
		fmt.Printf("%.2f + %.2f = %.2f\n",n,i+(s*float32(c)),n+(i+(s*float32(c))))
	}
	fmt.Printf("Tabuada de subtração : \n")
	for c:=0;c<k;c++{
		fmt.Printf("%.2f - %.2f = %.2f\n",n,i+(s*float32(c)),n-(i+(s*float32(c))))
	}
	fmt.Printf("Tabuada de multiplicação : \n")
	for c:=0;c<k;c++{
		fmt.Printf("%.2f x %.2f = %.2f\n",n,i+(s*float32(c)),n*(i+(s*float32(c))))
	}
	fmt.Printf("Tabuada de divisão : \n")
	for c:=0;c<k;c++{
		fmt.Printf("%.2f / %.2f = %.2f\n",n,i+(s*float32(c)),n/(i+(s*float32(c))))
	}
}
