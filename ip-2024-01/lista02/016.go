package main
import("fmt")

func cat(h int )[]int{
	var resp []int
	for c:=1;c<h;c++{
		for j:=1;j<h;j++{
			if h*h == (c*c)+(j*j){
				resp = append(resp, c)
				resp = append(resp, j)
				return resp
			}
		}
	}
	resp = append(resp, 0)
	resp = append(resp, 0)
	return resp
}

func main(){
	var n int
	fmt.Scan(&n)
	for c:=5;c<=n;c++{
		reg := cat(c)
		if reg[0]==0{
			continue
		}
		fmt.Printf("hipotenusa = %v, catetos %v e %v\n",c,reg[0],reg[1])
	}
}