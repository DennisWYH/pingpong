package util

func CheckIfSymbols(word string) bool {
	var result bool
	switch word {
	case ".":
		result = true
	case "。":
		result = true
	case ",":
		result = true
	case "，":
		result = true
	case "=":
		result = true
	case ")":
		result = true
	case "(":
		result = true
	case "&":
		result = true
	case "%":
		result = true
	default:
		result = false
	}
	return result
}
