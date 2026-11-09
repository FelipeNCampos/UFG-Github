package Models

func Catch_Gen(id_gen int) string {

		switch id_gen {
		case 1:
			return "Mulher cis"
		case 2:
			return "Homen cis"
		case 3:
			return "Mulher trans"
		case 4:
			return "Mulher trans"
		case 5:
			return "Não-binário"
		case 6:
			return "Gênero neutro"
		case 7:
			return "Outros"
		case 8:
			return "Prefere não informar"
		default:
			return ""
		}
}