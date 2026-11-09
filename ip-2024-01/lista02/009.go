package main
import("fmt")


func primo (n int)string{
	for c:=n-1;c>1;c--{
		if n%c==0{
			return "NAO PRIMO"
		}
	}
	return "PRIMO"
}

func main(){
	var n int 
	fmt.Scan(&n)
	fmt.Print(primo(n))
}