package main
import("fmt")

func main(){
    var n,temp,maior,menor int 
    var reg []int
    fmt.Scan(&n)
    for c:=0;c<n;c++{
        fmt.Scan(&temp)
		if c==0{
			menor = temp
		}else{
			if temp<menor{
				menor = temp
			}
		}
        reg = append(reg,temp)
		if temp>maior{
			maior = temp
		}
    }
	for c:=0;c<n;c++{
        fmt.Printf("%v ",reg[c])
    }
	fmt.Printf("\n")
    for c:=n-1;c>=0;c--{
        fmt.Printf("%v ",reg[c])
    }
	fmt.Printf("\n")
	fmt.Printf("%v\n%v",maior,menor)
}
