package main
import("fmt")

func main (){
	var x,y int
	fmt.Scan(&x)
	fmt.Scan(&y)
	for c:=1;c<=x;c++{
		for j:=1;j<=y;j++{
			if c>j{
				fmt.Printf("(%v,%v)",c,j)
				if j!=c-1 && j!=y{
					fmt.Print("-")
				}
			}
		}
		fmt.Print("\n")
	}
}