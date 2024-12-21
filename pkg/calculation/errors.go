package calculation

import "errors"

var (
	ErrInvalidExpression         = errors.New("invalid expression")
	ErrInvalidCharInExpression   = errors.New("invalid char in expression")
	ErrDivisionByZero            = errors.New("division by zero")
	ErrOpeningParenthesisMissing = errors.New("opening parenthesis missing")
	ErrClosingParenthesisMissing = errors.New("closing parenthesis missing")
)
