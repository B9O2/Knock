package components

func Ternary(condition bool, arg1, arg2 any) any {
	if condition {
		return arg1
	} else {
		return arg2
	}
}
