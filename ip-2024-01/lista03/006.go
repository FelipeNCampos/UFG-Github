package main
import("fmt")

func main(){
	var sum, temp,n int
	fmt.Scan(&n)
	for c:=0;c<n,c++{
		fmt.Scan(&temp)
		sum +=temp
	} 
	fmt.Printf("%v",sum)
}
