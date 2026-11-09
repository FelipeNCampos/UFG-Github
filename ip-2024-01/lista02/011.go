package main 
import(
	"fmt"
	"math"
	)

func rk(n float64,k int )float64{
	if k==0{
		return 1
	}else{
		return (rk(n,k-1)+(n/rk(n,k-1)))/2
	}
}

func main(){
	var n, e float64
	fmt.Scan(&n)
	fmt.Scan(&e)
	k:=1
	for {
		temp := math.Abs(n-rk(n,k)*rk(n,k))
		fmt.Printf("r : %.9f, erro: %.9f \n",rk(n,k),temp)
		if temp<e{
			break
		}
		k++
	}
}