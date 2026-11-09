package main 
import("fmt")

func main (){
	var temp,logica,i,r,e int 
	fmt.Scan(&i)
	fmt.Scan(&r)
	fmt.Scan(&e)
	for c:=0;c<e;c++{
		logica =i+(r*c)
		temp += logica
	}
	fmt.Printf("%v",temp)
}
