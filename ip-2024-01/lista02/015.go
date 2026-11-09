package main
import("fmt")

func inv(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func main(){
	var n int 
	var temp1,temp2 string
	fmt.Scan(&n)
	for c:=0;c<n;c++{
		fmt.Scan(&temp1)
		fmt.Scan(&temp2)
		if inv(temp1) > inv(temp2){
			fmt.Print(inv(temp1),"\n")
		}else{
			fmt.Print(inv(temp2),"\n")
		}
	}
}