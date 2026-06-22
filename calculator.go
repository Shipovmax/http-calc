package main

import (
	"fmt"
)

// calculate выполняет арифметику и возвращает ошибку вместо паники
func calculate(a float64, op string, b float64) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		// деление на ноль — ошибка, а не паника
		if b == 0 {
			return 0, fmt.Errorf("деление на ноль")
		}
		return a / b, nil
	default:
		// любой неизвестный оператор → сообщаем какой именно
		return 0, fmt.Errorf("неизвестный оператор: %s", op)
	}
}
