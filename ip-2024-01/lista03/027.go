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

func main(){
	var x,z,temp int 
	var reg1,reg2 []int
	fmt.Scan(&x)
	fmt.Scan(&z)
	for c:=0;c<x;c++{
		fmt.Scan(&temp)
		reg1 = append(reg1,temp)
	}
	for c:=0;c<z;c++{
		fmt.Scan(&temp)
		reg2 = append(reg2,temp)
	}
	reg := make([]int,len(reg1)+len(reg2))
	reg = sort(reg1,reg2)
}
