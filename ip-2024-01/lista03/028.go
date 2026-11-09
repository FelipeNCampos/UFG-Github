package main
import("fmt")

func main(){
	var x,z,temp,flag int
	var regx,regz,reg,repetidos []int
	for{
		fmt.Scan(&temp)
		if temp>0 || temp<=1000{
			x = temp
			break
		}
	}
	for{
		fmt.Scan(&temp)
		if temp>0 || temp<=1000{
			z=temp
			break
		}
		
	}
	for c:=0;c<x;c++{
		fmt.Scan(&temp)
		regx = append(regx,temp)
	}
	for c:=0;c<z;c++{
		fmt.Scan(&temp)
		regz = append(regz,temp)
	}
	for c:= range regx{
		flag=0
		for j:= range reg{
			if reg[j]==regx[c]{
				flag = 1
				repetidos = append(repetidos,regx[c])
			}
		}
		if flag==0{
			reg = append(reg,regx[c])
		}
	}
	for c:= range regz{
		flag=0
		for j:= range reg{
			if reg[j]==regz[c]{
				flag = 1
				repetidos = append(repetidos,regz[c])
			}
		}
		if flag==0{
			reg = append(reg,regz[c])
		}
	}
	fmt.Print("uniao : :",reg,"\n","repetidos:",repetidos)
	
}
