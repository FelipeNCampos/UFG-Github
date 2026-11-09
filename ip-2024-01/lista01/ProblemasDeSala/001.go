package main

import (
"fmt"
"strconv"
)

func main() {
  var arr []int
  for c := 0; c < 3; c++ {
    var input string
    fmt.Println("Digite o %d* número: ", c+1)
    fmt.Scanln(&input)
    num, err := strconv.Atoi(input)
  if err != nil {
    fmt.Println("Erro ao converter a entrada para número inteiro.")
    return
  }
  arr = append(arr, num)
}

for c := 0; c < 3; c++ {
  for j := c + 1; j < 3; j++ {
    if arr[c] > arr[j] {
      temp := arr[c]
      arr[c] = arr[j]
      arr[j] = temp
      }
    }
  }

fmt.Println(arr)
}