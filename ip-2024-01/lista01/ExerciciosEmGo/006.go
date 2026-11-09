package main

import ("fmt")

func main (){
	var temp int
	var reg float32 
	fmt.Scan(&temp)
	for c:=0;c<temp;c++{
		fmt.Scan(&reg)
		fmt.Printf("%.2f FAHRENHEIT EQUIVALE A %.2f CELSIUS\n",reg,((reg-32)*5)/9)
	}
}