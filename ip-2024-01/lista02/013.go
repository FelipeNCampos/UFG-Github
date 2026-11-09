package main
import(
	"fmt"
	)

func main (){
	var n,sum float64
	fmt.Scan(&n)
	for c:=1;c<64;c++{
		if c%2==1{
			sum += n*2
		}else{
			sum += n
		}
	}
	fmt.Printf("%v",sum)
}