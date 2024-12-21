package calculation

import (
	"strconv"
	"strings"
	"unicode"
)

func toRPN(expression string) ([]string, error) {
	var output []string
	var operators []rune
	var numBuffer strings.Builder

	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	isUnary := true

	for _, char := range expression {
		if unicode.IsSpace(char) {
			continue
		}

		if unicode.IsDigit(char) || char == '.' {
			numBuffer.WriteRune(char)
			isUnary = false
		} else {
			if numBuffer.Len() > 0 {
				output = append(output, numBuffer.String())
				numBuffer.Reset()
			}

			switch char {
			case '(':
				operators = append(operators, char)
				isUnary = true
			case ')':
				for len(operators) > 0 && operators[len(operators)-1] != '(' {
					output = append(output, string(operators[len(operators)-1]))
					operators = operators[:len(operators)-1]
				}
				if len(operators) == 0 {
					return nil, ErrOpeningParenthesisMissing
				}
				operators = operators[:len(operators)-1]
				isUnary = false
			case '+', '-':
				if isUnary {
					if char == '-' {
						output = append(output, "0")
						operators = append(operators, '-')
					} else if char == '+' {
						continue
					}
					isUnary = false
					continue
				}
				for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[char] {
					output = append(output, string(operators[len(operators)-1]))
					operators = operators[:len(operators)-1]
				}
				operators = append(operators, char)
				isUnary = true
			case '*', '/':
				for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[char] {
					output = append(output, string(operators[len(operators)-1]))
					operators = operators[:len(operators)-1]
				}
				operators = append(operators, char)
				isUnary = true
			default:
				return nil, ErrInvalidCharInExpression
			}
		}
	}

	if numBuffer.Len() > 0 {
		output = append(output, numBuffer.String())
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == '(' {
			return nil, ErrClosingParenthesisMissing
		}
		output = append(output, string(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

func evalRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, ErrDivisionByZero
				}
				stack = append(stack, a/b)
			}
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}

func Calc(expression string) (float64, error) {
	rpn, err := toRPN(expression)
	if err != nil {
		return 0, err
	}

	return evalRPN(rpn)
}
