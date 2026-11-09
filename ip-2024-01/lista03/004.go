package main
import("fmt")

func main(){
    var n,temp int 
    var reg []int
    fmt.Scan(&n)
    for c:=0;c<n;c++{
        fmt.Scan(&temp)
        reg = append(reg,temp)
    }
    for c:=0;c<n;c++{
        fmt.Printf("%v ",reg[c])
    }
}
