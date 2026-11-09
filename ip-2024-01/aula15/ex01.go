package main

import("fmt")

func exp(num int,e int )int{
	if e==1{
		return num
	}else{
		return num*exp(num,e-1)
	}
}

func main(){
	var num, e int 
	fmt.Scan(&num)
	fmt.Scan(&e)
	fmt.Print(exp(num,e))
}
