package main
import "fmt"

func main() {
    var m,d float32
    fmt.Print("Digite o preço do produto : ")
    fmt.Scan(&m)
    fmt.Print("Digite o desconto no produto : ")
    fmt.Scan(&d)
    fmt.Print("O valor do produto é :")
    fmt.Println(m*(1-(d/100)))
}
