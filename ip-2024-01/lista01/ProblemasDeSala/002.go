package main

import(
    "fmt"
    "math"
    "strconv"
    )

func main(){
  var arr []int
  for c := 0; c < 3; c++ {
    var input string
    fmt.Printf("Digite o %d* número:\n", c+1)
    fmt.Scanln(&input)
    num, err := strconv.Atoi(input)
    if err != nil {
      fmt.Println("Erro ao converter a entrada para número inteiro.")
      return
    }
    arr = append(arr, num)
    }
    max := math.MinInt64
    for _, v := range arr {
      if v > max {
        max = v
      }
    }
  min := math.MaxInt64
  for _, v := range arr {
    if v < min {
      min = v
    }
  }
  fmt.Println("Menor = ",min)
  fmt.Println("Maior = ",max)
  }