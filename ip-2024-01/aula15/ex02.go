package main

import("fmt")

func somaarr(arr []int,id int)int{
	if id ==0{
		return 0
	}else{
		return arr[tam-1] + somaarr(arr,id-1)
	}
}


func main(){
	var n,temp int
	var arr []int
	fmt.Print("Digite o tamanho do array :")
	fmt.Scan(&n)
	for c:=0;c<n;c++{
		fmt.Printf("Digite o %v* valor do array :",c+1)
		fmt.Scan(&temp)
		arr = append(arr,temp)
	}
	fmt.Print(somaarr(arr),len(arr))
}
