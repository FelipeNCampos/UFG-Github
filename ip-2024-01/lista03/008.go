package main
import("fmt")


func conv(n int){
    var cont int
	var reg []int
	for {
		reg = append(reg,n%2)
		n = n/2
		if n/2==0{
		    cont++
		}
		if cont==2{
		    break
		}
	}
		for c:= range reg{
			fmt.Printf("%v",reg[len(reg)-(c+1)])
		}
		fmt.Printf("\n")
}

func main(){
	for{
		var temp int
		fmt.Scan(&temp)
		conv(temp)
	}
}
