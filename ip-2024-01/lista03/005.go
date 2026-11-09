package main
import("fmt")

func main(){
	var temp,num,id int
	for{
	    max := 0
		fmt.Scan(&temp)
		if temp==0{
			break
		}
		for c:=0;c<temp;c++{
			fmt.Scan(&num)
			if num>max{
				max = num
				id = c
			}
		}
		fmt.Printf("%v %v \n",id , max)
        
	}
}
