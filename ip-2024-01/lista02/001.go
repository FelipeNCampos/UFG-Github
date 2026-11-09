package main
import("fmt")

func main(){
	for {
		var matricula int 
		fmt.Scan(&matricula)
		if matricula == -1{
			return
		}
		var temp,mp,ml,nt float32
		for c:=0;c<8;c++{
			fmt.Scan(&temp)
			mp +=temp
		}
		mp /=8
		for c:=0;c<5;c++{
			fmt.Scan(&temp)
			ml +=temp
		}
		ml /=5
		fmt.Scan(&temp)
		nt += temp
		var pres float32
		fmt.Scan(&pres)
		if (0.7*mp + 0.15*ml +0.15*nt>6){
			if pres>96{
				fmt.Printf("Matricula : %v, Nota final: %.2f, Situação final:Aprovado \n",matricula,float32(0.7*mp + 0.15*ml +0.15*nt))
			}else{
				fmt.Printf("Matricula : %v, Nota final: %.2f, Situação final:Reprovado por frequencia \n ",matricula,float32(0.7*mp + 0.15*ml +0.15*nt))
			}
		}else{
			if pres>96{
				fmt.Printf("Matricula : %v, Nota final: %.2f, Situação final:Reprovado por nota \n",matricula,float32(0.7*mp + 0.15*ml +0.15*nt))
			}else{
				fmt.Printf("Matricula : %v, Nota final: %.2f, Situação final:Reprovado por nota e frequencia \n",matricula,float32(0.7*mp + 0.15*ml +0.15*nt))
			}
		}
	}
}
