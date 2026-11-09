package main

import"fmt"

func main() {
    var N int
    
    fmt.Scan(&N)
    
    for c:=0 ; c < N; c++ {
        var reg []int
        var num,temp int 
        fmt.Scan(&num)
        for j :=0 ;j<num;j++{
            fmt.Scan(&temp)
            reg = append(reg,temp)
        }
        tam := len(reg)/2
        for j:=0;j<=tam;j++{
            if reg[j]<reg[len(reg)-1-j]{
                reg[j],reg [len(reg)-1-j] = reg [len(reg)-1-j] , reg[j]
            }
        }
        fmt.Print(reg)
    }
}