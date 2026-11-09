package main
import("fmt")

var num int
func fac(n int)int{
	if n==1{
		return n 
	}else{
		return fac(n-1)*n
	}
}


func main(){
	fmt.Scan(&num)
	fmt.Printf("%v!=%v\n",num,fac(num))
}
