package main
import ("fmt")

func main(){
	var n, temp,sena,quadra,quina int
	var reg []int
	for c:=0;c<6;c++{
		fmt.Scan(&temp)
		reg = append(reg, temp)
	}
	fmt.Scan(&n)
	for c:=0;c<n;c++{
		cont:=0
		for z :=0;z<6;z++{
			fmt.Scan(&temp)
			for j:= range reg{
				if temp==reg[j]{
					cont++
				}
			}
		}
		if cont==6{
			sena++
		}
		if cont==5{
			quina++
		}
		if cont==4{
			quadra++
		}
	}
	if sena>0{
		fmt.Printf("Houve %v acertador(es) da sena\n",sena)
	}else{
		fmt.Printf("Nao houve acertador da sena\n")
	}
	if quina>0{
		fmt.Printf("Houve %v acertador(es) da quina\n",quina)
	}else{
		fmt.Printf("Nao houve acertador da quina\n")
	}
	if quadra>0{
		fmt.Printf("Houve %v acertador(es) da quadra\n",quadra)
	}else{
		fmt.Printf("Nao houve acertador da quadra\n")
	}
}