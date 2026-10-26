package domain

// escalaFibonacci es la escala admitida para el valor de negocio y los Story Points.
var escalaFibonacci = map[int]struct{}{
	1:  {},
	2:  {},
	3:  {},
	5:  {},
	8:  {},
	13: {},
	21: {},
}

// ErrEstimacionInvalida indica unos Story Points fuera de la escala Fibonacci.
var ErrEstimacionInvalida = ValidationError{Campo: "estimacion_sp", Mensaje: "los Story Points deben ser uno de 1, 2, 3, 5, 8, 13, 21"}

// EsValorNegocioValido indica si v pertenece a la escala Fibonacci permitida.
func EsValorNegocioValido(v int) bool {
	_, ok := escalaFibonacci[v]
	return ok
}

// EsEstimacionValida indica si los Story Points pertenecen a la escala Fibonacci permitida.
func EsEstimacionValida(v int) bool {
	_, ok := escalaFibonacci[v]
	return ok
}
