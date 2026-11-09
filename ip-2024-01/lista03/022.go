package main
import ("fmt")

func main(){
    for{
        var t,f int
        var ent,resp string
        fmt.Scan(&t)
        fmt.Scan(&f)
        if t==0{
            break
        }
        fmt.Scan(&ent)
        for c:=t-1;c>=0;c--{
            cont:=0
            if len(resp)<f{
                resp = string(ent[c:])
            }else{
                if string(ent[c])>string(resp[cont]){
                    resp = string(ent[c])+string(resp[cont+1:])
                }
            }
        }
        fmt.Print(resp,"\n")
        
    }
}