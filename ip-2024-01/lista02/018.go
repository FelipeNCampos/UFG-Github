package main 
import ("fmt")

func main(){
	var a,b,c int 
	fmt.Scan(&a)
	fmt.Scan(&b)
	fmt.Scan(&c)
	cont := 2
	mmc := 1
	for{
		if(a==1&&b==1&&c==1){
			break
		}
		if(a%cont!=0 && b%cont!=0 && c%cont!=0){
			cont++
		}else {
			mmc *= cont
			fmt.Printf("%v %v %v | %v\n",a,b,c,cont)		
			if a%cont==0{
				a = a/cont
			}
			if b%cont==0{
				b = b/cont
			}
			if c%cont==0{
				c = c/cont
			}
		}
	}
	fmt.Printf("MMC : %v\n",mmc)
}
