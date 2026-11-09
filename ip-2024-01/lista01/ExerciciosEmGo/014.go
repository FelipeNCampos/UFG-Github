package main 
import (
	"fmt"
	)

func main(){
	var h,a float32
	fmt.Scan(&h)
	fmt.Scan(&a)
	fmt.Printf("O VOLUME DA PIRAMIDE E = %.2f METROS CUBICOS.",((3*(a*a)*1.73205)/2)*h/3)
}
