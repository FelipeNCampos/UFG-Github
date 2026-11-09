package main
import(
	"fmt"
	"strconv"
)

func main(){
	var resultado string
	for c:=0;c<3;c++{
		var temp int
		fmt.Scan(&temp)
		if len(strconv.Itoa(temp))>1{
			fmt.Print("DIGITO INVALIDO")
			return
		}
		resultado = resultado + strconv.Itoa(temp)
	}
	intx , _ := strconv.ParseInt(resultado,10,64)
	fmt.Print(resultado," ",intx*intx)
}