package Database


var Atual int
var att bool
var busca string

func SetSession(id int) {
	Atual = id
	att = true
}

func GetSession() int {
	return Atual
}

func LogOff() {
	att = false
}

func CheckSesison() bool {
	return att
}


func SetFind(cpf string){
	busca = cpf
}

func GetFind()string{
	return busca
}