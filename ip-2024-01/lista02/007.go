package main
import ("fmt")

func media(numbers []int) float32 {
	var ind,sum float32
    for _, num := range numbers {
        sum += float32(num)
		ind++
    }
    return sum/ind
}

func main (){
	var pares []int
	var impares []int
	var temp int
	for {
		fmt.Scan(&temp)
		if temp==0{
			break
		}
		if temp%2==0{
			pares = append(pares,temp)
		}else{
			impares = append(impares,temp)
		}
	}
	var mp,mi float32
	mp = media(pares)
	fmt.Printf("media par : %.2f\n",mp)
	mi = media(impares)
	fmt.Printf("media impar : %.2f\n",mi)

}
