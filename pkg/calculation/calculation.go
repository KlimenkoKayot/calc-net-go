package calculation

import (
	"github.com/klimenkokayot/calc-api-go/pkg/rpn"
)

func Calc(expression string) (float64, error) {
	values, err := rpn.ExpressionToRPN(expression)
	if err != nil {
		return 0.0, err
	}
	val := []float64{}
	for i := 0; i < len(values); i++ {
		if num, ok := values[i].(float64); ok {
			val = append(val, num)
		} else {
			if len(val) < 2 {
				return 0.0, ErrInvalidExpression
			}

			b := val[len(val)-1]
			val = val[:len(val)-1]
			a := val[len(val)-1]
			val = val[:len(val)-1]

			switch values[i].(string) {
			case "+":
				val = append(val, a+b)
			case "-":
				val = append(val, a-b)
			case "*":
				val = append(val, a*b)
			case "/":
				if b == 0 {
					return 0.0, ErrDevisionByZero
				}
				val = append(val, a/b)
			}
		}
	}
	if len(val) > 1 {
		return 0.0, ErrInvalidExpression
	}
	if len(val) == 0 {
		return 0.0, ErrEmptyExpression
	}
	return val[0], nil
}
