package main

import(
	"fmt"
)


func invarr(arr []int,id int){
	if id==0{
		print(arr[id]," ")
		return 
	}else{
		print(arr[id]," ")
		invarr(arr,id-1)
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
	invarr(arr,len(arr)-1)
}
