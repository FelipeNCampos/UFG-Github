package main
import("fmt")

func main (){
	var n,temp,resp int 
	var reg []int
	fmt.Scan(&n)
	for c:=0;c<n;c++{
		fmt.Scan(&temp)
		reg = append(reg, temp)
	}
	for c:=0;c<n-1;c++{
		cont := 1
		for j:=c+1;j<n;j++{
			if reg[c]==reg[j]{
				cont++
			}else{
				break
			}
		}
		if cont>=resp{
			resp=cont
		}
	}
	fmt.Printf("A maior subsequencia de elementos iguais possui %v elementos .",resp)
}
