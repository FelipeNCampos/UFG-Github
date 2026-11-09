package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
)

func main(){
	var char rune
	var reg string
	var err error
	fmt.Println("input:")
    reader := bufio.NewReader(os.Stdin)
	for{
		char,_,err = reader.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if string(char)=="#"{
			break
		}else{
			reg = string(char)+reg
		}
	}
	fmt.Print("output:\n", reg)
}