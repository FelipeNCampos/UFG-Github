package main
import("fmt")

func main(){
	var n,temp,lk,cont,resp int 
	fmt.Scan(&n)
	var reg []int
	for c:=0;c<n;c++{
		fmt.Scan(&temp)
		reg = append(reg, temp)
	}
	for c:=0;c<n-1;c++{
		lk = reg[c]
		for j:=c+1;j<n;j++{
			if reg[j]>lk{
				lk = reg[j]
				cont++
			}else{
				if cont>=resp{
					resp = cont
				}
				break
			}
			
			}
			cont=0
		}
		fmt.Printf("O comprimento do segmento crescente maximo e : %v",resp)
	}
	

