package main
import("fmt")

func findM(list []float64) (int, float64, int, float64) {
	var min, max float64 = list[0], list[0]
	var idx1, idx2 int = 0, 0
	for i := range list {
		if list[i] < min {
			min = list[i]
			idx1 = i
		}
		if list[i] > max {
			max = list[i]
			idx2 = i
		}
	}
	return idx1, min, idx2, max
}

func IndInputs(params ...int) []float64 {
	var n int
	if len(params) == 0 {
		fmt.Printf("Digite o número de %v: ", params[0])
		fmt.Scan(&n)
	} else {
		n = params[0]
	}

	v := []float64{}
	for i := 0; i < n; i++ {
		var x int
		fmt.Scan(&x)
		v = append(v, float64(x))
	}
	return v
}

func Sum(list []float64) float64 {
	var s float64
	for i := range list {
		s += list[i]
	}
	return s
}

func sortir(list []float64) []float64 {
	var sorted, cp = make([]float64, 0, len(list)), list[:]
	for i := 0; i < len(list); i++ {
		_, min, _, _ := findM(cp)
		sorted = append(sorted, min)
		for j := range cp {
			if cp[j] == min && len(cp) > 1 {
				cp[j] = cp[len(cp)-1]
				cp = cp[:len(cp)-1]
				break
			}
		}
	}
	return sorted
}

func main() {
	var n int
	fmt.Print("Digite o número de casos de teste: ")
	fmt.Scan(&n)
	for k := 0; k < n; k++ {
		fmt.Print("Digite o número de cada anão: \n")
		v := sortir(IndInputs(9))
	out:
		for i := range v {
			for j := range v {
				if i != j && Sum(v)-v[i]-v[j] == 100 {
					fmt.Print("Os anões são os de número:\n")
					for anao := range v {
						if anao != i && anao != j {
							fmt.Println(v[anao])
						}
					}
					fmt.Println("\n")
					break out
				}
			}
		}
	}
}