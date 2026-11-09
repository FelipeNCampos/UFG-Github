package main 
import("fmt")

func main(){
	var temp int
	fmt.Scan(&temp)
	for c:=2;c<=temp;c+=2{
		fmt.Printf("%v^2 = %v\n",c,c*c)
	}
}
