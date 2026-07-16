package main

import (
	"fmt"
)

// calculate performs the arithmetic operation op on a and b.
// It returns an error instead of panicking on division by zero
// or an unsupported operator.
func calculate(a float64, op string, b float64) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		// division by zero is a domain error, not a panic
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	default:
		// report the unsupported operator explicitly
		return 0, fmt.Errorf("unknown operator: %s", op)
	}
}
