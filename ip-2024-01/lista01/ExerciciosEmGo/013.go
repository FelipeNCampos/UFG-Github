package main
import("fmt")

func main (){
	var temp float32
	fmt.Scan(&temp)
	if temp>=9{
		fmt.Printf("NOTA = %.1f , CONCEITO = A ", temp)
	}else{
		if temp>=7.5{
			fmt.Printf("NOTA = %.1f , CONCEITO = B ", temp)
		}else{
			if temp>=6{
				fmt.Printf("NOTA = %.1f , CONCEITO = C ", temp)
			}else{
				fmt.Printf("NOTA = %.1f , CONCEITO = D ", temp)
			}
		}
	}
}
