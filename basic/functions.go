package main

type UnaryFunction func(int) int

func sqr(a int) int {
	return a * a
}

func revert(a int) int {
	return -1 * a
}

func combine(first UnaryFunction, second UnaryFunction) UnaryFunction {
	return func(a int) int {
		return second(first(a))
	}
}

func Map(slice []int, fun UnaryFunction) {
	for i, val := range slice {
		slice[i] = fun(val)
	}
}
