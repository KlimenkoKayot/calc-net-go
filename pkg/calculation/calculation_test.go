// TODO
// Нужно сделать 100% покрытие тестами файла calculation_test.go

package calculation

import (
	"testing"
)

type Test struct {
	name       string
	expression string
	expected   float64
	err        error
}

func TestCalc(t *testing.T) {
	cases := []Test{
		{
			name:       "Default expression",
			expression: "2+2*2",
			expected:   6,
		},
		{
			name:       "Default expression with ()",
			expression: "(2+2)*2",
			expected:   8,
		},
		{
			name:       "Big default expression",
			expression: "2+2*2+3*(5+1)*2/2/5/1+3/4",
			expected:   10.35,
		},
		{
			name:       "Expression with spaces",
			expression: "(   (2  )  )",
			expected:   2,
		},
		{
			name:       "Hard expression",
			expression: "((2-2)+1)/2",
			expected:   0.5,
		},
		{
			name:       "ERROR only ()",
			expression: "(()",
			err:        ErrInvalidExpression,
		},
		{
			name:       "ERROR invalid symbol",
			expression: "g",
			err:        ErrInvalidSymbolExpression,
		},
		{
			name:       "ERROR devision by zero",
			expression: "2/0",
			err:        ErrDevisionByZero,
		},
		{
			name:       "ERROR empty expression",
			expression: "",
			err:        ErrEmptyExpression,
		},
		{
			name:       "ERROR invalid expression #1",
			expression: "+)",
			err:        ErrInvalidExpression,
		},
		{
			name:       "ERROR invalid expression #2",
			expression: "-+-2+",
			err:        ErrInvalidExpression,
		},
		{
			name:       "ERROR invalid expression #3",
			expression: "2+2(2)(2)(20)",
			err:        ErrInvalidExpression,
		},
		{
			name:       "ERROR invalid operation",
			expression: "2?3",
			err:        ErrInvalidSymbolExpression,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Calc(tc.expression)
			if err != tc.err {
				t.Fatalf("Calc(\"%s\"), error: %v, expected: %v", tc.expression, err, tc.err)
			}
			if got != tc.expected {
				t.Fatalf("Calc(\"%s\"), got: %v, expected: %v", tc.expression, got, tc.expected)
			}
		})
	}
}
