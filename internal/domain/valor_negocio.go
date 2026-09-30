package domain

// valoresNegocioValidos es la escala Fibonacci admitida para el valor de negocio.
var valoresNegocioValidos = map[int]struct{}{
	1:  {},
	2:  {},
	3:  {},
	5:  {},
	8:  {},
	13: {},
	21: {},
}

// EsValorNegocioValido indica si v pertenece a la escala Fibonacci permitida.
func EsValorNegocioValido(v int) bool {
	_, ok := valoresNegocioValidos[v]
	return ok
}
