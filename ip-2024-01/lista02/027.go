package main
import ("fmt")

func main(){
	var n int
	for{
		fmt.Scan(&n)
		if n < 1{
			fmt.Printf("Fatoracao nao e possive para o numero %v !\n",n)
		}else{
			break
		}
	}
		if n==1{
			fmt.Printf("1 = 1")
		}else{
			fmt.Printf("%v = ",n)
			fat := 2
			for{
				if n%fat!=0{
					fat++
					continue
				}
				n = n/fat
				fmt.Printf("%v",fat)
				if n ==1{
					break
				}
				if n != 1{
					fmt.Print(" x ")
				}
			}
		
	}
}