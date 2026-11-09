package main
import ("fmt"
"bufio"
"math"
"os")

func BToi(stat bool) int {
	if stat {
		return 1
	} else {
		return 0
	}
}

func main(){
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	vvv := make([]int, 11)
	sep := false
	for i := 0; i < len(str); i++ {
		Vogais := [][]string{{"a", "A"}, {"e", "E"}, {"i", "I"}, {"o", "O"}, {"u", "U"}, {";", ";"}}
		for j := 0; j < len(Vogais); j++ {
			if string(str[i]) == Vogais[j][0] || string(str[i]) == Vogais[j][1] {
				if j == 5 {
					sep = true
					vvv[10] += 1
				} else {
					vvv[j+[]int{0, 5}[BToi(sep)]] += 1
				}
			}
		}
	}
	if vvv[10] != 1 {
		fmt.Print(vvv[10], "Formato inválido.\n")
		return
	}
	var distance float64
	for i := 0; i < 5; i++ {
		distance += math.Pow(float64(vvv[i]-vvv[i+5]), 2)
	}
	distance = math.Sqrt(distance)
	fmt.Printf("%v\n%v\n%.2f\n", vvv[:5], vvv[5:10], distance)
}