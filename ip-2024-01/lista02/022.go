package main
import(
	"fmt"
	"math"
	)

func main(){
	var n float64
	var num,den int
	cont := 2.0
	fmt.Scan(&n)
	for {
		if math.Mod(n*cont,2)==1 || math.Mod(n*cont,2)==0{
			den = int(cont)
			num = int(n*cont)
			break
		}
		cont++
	}
	fmt.Printf("%v/%v",num,den)
}