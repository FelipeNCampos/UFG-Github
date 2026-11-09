package main
import("fmt")

func main (){
	for{
		var matricula int
		var horas,vph float32
		fmt.Print("input : ")
		fmt.Scan(&matricula)
		if matricula==0{
			return 
		}
		fmt.Scan(&horas)
		fmt.Scan(&vph)
		fmt.Print("output : ")
		fmt.Printf("%v %.2f\n",matricula,horas*vph)
		fmt.Print("------------------------------\n")
	}
}