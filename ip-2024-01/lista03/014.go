package main
import("fmt")

func sort(Ll,Rl []int)[]int{
	var i,j,k int
	var resp []int
	temp := make([]int,len(Ll)+len(Rl))
	for{
		if i>=len(Ll) || j>=len(Rl){
			break
		}
		if Ll[i]<Rl[j]{
			temp[k]=Ll[i]
			k++
			i++
		}else{
			temp[k]=Rl[j]
			k++
			j++
		}
	}
		for{
			if i>=len(Ll){
			break
			}
			temp[k]=Ll[i]
			i++
			k++
		}
		for{
			if j>=len(Rl){
				break
			}
			temp[k]=Rl[j]
			j++
			k++
		}
		for c:= range temp{
			resp = append(resp,temp[c])
		}
		return resp
	}


func mergesort(lista []int)[]int{
	n := len(lista)
	if n<2{
		return lista
	}
	mid := n/2
	var Llista []int
	var Rlista []int
	for c:=0;c<mid;c++{
		Llista = append(Llista, lista[c])
	}
	for c:=mid;c<n;c++{
		Rlista = append(Rlista, lista[c])
	}
	Llista = mergesort(Llista)
	Rlista = mergesort(Rlista)
	lista = sort(Llista,Rlista)
	return lista
}

func main(){
	var x,temp int
	var reg []int
	var resposta float64
	fmt.Scan(&x)
	for c:=0;c<x;c++{
		fmt.Scan(&temp)
		reg = append(reg, temp)
	}
	reg = mergesort(reg)
	if len(reg)%2==1{
		resposta = float64(reg[len(reg)/2])
	}else{
		resposta = float64(reg[len(reg)/2-1]+reg[len(reg)/2])/2
	}
	fmt.Printf("%.2f",resposta)
}