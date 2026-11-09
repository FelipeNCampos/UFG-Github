package main
import ("fmt")


func main(){
	var temp int
	fmt.Scan(&temp)
	for c:=0;c<temp;c++{
		var total float64
		fmt.Scan(&total)
		var popular float64
		fmt.Scan(&popular)
		var geral float64
		fmt.Scan(&geral)
		var arquibancada float64
		fmt.Scan(&arquibancada)
		var cadeiras float64
		fmt.Scan(&cadeiras)
		renda := total*popular/100 + (total*geral/100)*5 + (total*arquibancada/100)*10 + (total*cadeiras/100)*20
		fmt.Print("A RENDA DO JOGO N.",c+1,"E = ",renda,"\n")
		
	}

}