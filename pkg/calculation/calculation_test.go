package calculation

import (
	"testing"
)

func TestCalc(t *testing.T) {
	testCases := []struct {
		expression string
		expected   float64
	}{
		{"-2", -2},
		{"2+2", 2 + 2},
		{"6+ -2", 6 + -2},
		{"10*(12/6*-2)", 10 * 12 / 6 * -2},
		{"(3+3) *2", (3 + 3) * 2},
		{"7+ -3 + 9 - -4", 7 + -3 + 9 - -4},
		{"(5+4+3)/(2+1) * 9 / 3", (5 + 4 + 3) / (2 + 1) * 9 / 3},
		{"(48/6) * 5 /((3+1) * (2+3)) -3", (48/6)*5/((3+1)*(2+3)) - 3},
		{"-2+3", -2 + 3},
		{"4+4*5", 4 + 4*5},
		{"((8+2) / (3+2) * 6) / 9 * (30 - ((5+10)*2)) -2", ((8+2)/(3+2)*6)/9*(30-((5+10)*2)) - 2},
		{"-50+50", -50 + 50},
		{"12*-2", 12 * -2},
		{"(2+3*(12) + 9)", (2 + 3*(12) + 9)},
		{"((2+5) * (2+3) +12) *3", ((2+5)*(2+3) + 12) * 3},
		{"2+2+2+2", 2 + 2 + 2 + 2},
		{"-(6+7)", -(6 + 7)},
		{"3+6/3*4", 3 + 6/3*4},
		{"10*-2", 10 * -2},
		{"(2)", 2},
		{"((3+4)*(5*(6+3) - 48 / (2+4) * (1+2)) - (8-3)) + (9 * (4-3 * (2+2)))", ((3+4)*(5*(6+3)-48/(2+4)*(1+2)) - (8 - 3)) + (9 * (4 - 3*(2+2)))},
		{"+(3+9)", +(3 + 9)},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expression, func(t *testing.T) {
			result, err := Calc(testCase.expression)
			if err != nil {
				t.Errorf("Calculation of expression %s failed with error: %v", testCase.expression, err)
			} else if result != testCase.expected {
				t.Errorf("Calculation of expression %s = %v, but wanted %v", testCase.expression, result, testCase.expected)
			}
		})
	}
}

func TestCalcWithErrors(t *testing.T) {
	testCases := []string{
		"2+4-3)",
		"22^37",
		"7/0",
		"*42",
		"4+(5-1",
		"cucumber",
		"7*(12+6*(27+3*(4+9) + 4 * (11+2) + 5 )",
		"5*(15+2",
		"(((((8))))",
		"-+",
		"'25",
		"61k+5t",
		"собака",
		"-",
		"()",
		"",
		"7**5",
	}

	for _, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			_, err := Calc(testCase)
			if err == nil {
				t.Errorf("The error of calculation of expression %s is not nil", testCase)
			}
		})
	}
}
