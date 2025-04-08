package casing

type ccase int

const (
	undefined ccase = iota
	upper
	lower
)

func charCase(c byte) ccase {
	switch {
	case c >= 'A' && c <= 'Z':
		return upper
	case c >= 'a' && c <= 'z':
		return lower
	default:
		return undefined
	}
}

func low(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c - 'A' + 'a'
	}
	return c
}

// CamelToDots converts any input string in CamelCase or dromedaryCase to dot.case
func CamelToDots(in string) string {
	return camelToCharSeparator(in, '.')
}

// CamelToSnake converts any input string in CamelCase or dromedaryCase to snake_case
func CamelToSnake(in string) string {
	return camelToCharSeparator(in, '_')
}

func camelToCharSeparator(in string, separator byte) string {
	res := make([]byte, 0, 10*len(in)/8)
	current := undefined
	for _, c := range []byte(in) {
		prev := current
		current = charCase(c)
		if current == upper && prev == lower {
			res = append(res, separator)
		}
		res = append(res, low(c))
	}
	return string(res)
}
