package calc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/nonrep/go-homework-1-calc/stack"
)

const (
	numberToken    = 'N' // Функция tokenize преобразует формулу в набор токенов, где все числа заменяются токеном N.
	openingBracket = '('
	closingBracket = ')'
	space          = ' '
	dot            = '.'

	plusOperator           = '+'
	minusOperator          = '-'
	multiplicationOperator = '*'
	divisionOperator       = '/'
)

var priority = map[rune]int{
	openingBracket:         0,
	closingBracket:         1,
	plusOperator:           7,
	minusOperator:          7,
	divisionOperator:       8,
	multiplicationOperator: 8,
}

var (
	errBracketCombination = errors.New("wrong combination of brackets")
	errInvalidFormula     = errors.New("the formula contains an invalid combination")
)

// operation выполняет арифметической операции.
func operation(left, right float64, operation rune) (result float64, err error) {
	switch operation {
	case plusOperator:
		result = left + right
	case minusOperator:
		result = left - right
	case multiplicationOperator:
		result = left * right
	case divisionOperator:
		if right == 0 {
			return result, errors.New("division by zero")
		}
		result = left / right
	}
	return result, nil
}

// isValidFormula проверяет наличие недопустимых символов в формуле.
func isValidFormula(formula string) error {
	validRunes := "0123456789 +-*/()."
	for _, rune := range formula {
		if !(strings.ContainsRune(validRunes, rune)) {
			return errors.New("formula is invalid")
		}
	}
	return nil
}

// addUnaryMinus преобразует формулу для поддержки унарного минуса.
func addUnaryMinus(formula string) string {
	var prev rune
	for i, char := range formula {
		if char == minusOperator && (prev == 0 || prev == space || prev == openingBracket) {
			formula = formula[:i] + "0" + formula[i:] // Добавляем 0 перед унарным минусом.
		}
		if char != space {
			prev = char
		}
	}
	return formula
}

// tokenize преобразует строку в слайс токенов, все числа заменяются на токен N и помещаются в слайс numbers.
func tokenize(formula string) (tokens []rune, numbers []float64, err error) {
	var numberString string

	for i := 0; i < len(formula); i++ {
		char := rune(formula[i])

		if unicode.IsDigit(char) || char == dot {
			numberString += string(char)
		} else {
			if numberString != "" {
				number, err := strconv.ParseFloat(numberString, 64)
				if err != nil {
					return nil, nil, errors.New("invalid number format")
				}
				tokens = append(tokens, numberToken)
				numbers = append(numbers, number)
				numberString = ""
			}

			if char == plusOperator || char == minusOperator || char == multiplicationOperator ||
				char == divisionOperator || char == openingBracket || char == closingBracket {
				tokens = append(tokens, char)
			}
		}
	}
	if numberString != "" {
		number, err := strconv.ParseFloat(numberString, 64)
		if err != nil {
			return nil, nil, errors.New("invalid number format")
		}
		tokens = append(tokens, numberToken)
		numbers = append(numbers, number)
	}

	if len(tokens) == 0 {
		return nil, nil, errors.New("the formula is empty or invalid")
	}

	return tokens, numbers, nil
}

// infixToPostfix определяет приоритет операций с помощью обратной польской записи, на выходе формула в виде слайса токенов в постфиксной записи.
func infixToPostfix(tokens []rune) ([]rune, error) {
	// В infix записи между операндами должен быть оператор.
	if strings.Contains(string(tokens), "NN") {
		return []rune{}, errInvalidFormula
	}

	var polishEntry []rune

	stack := stack.New[rune]()

	for _, token := range tokens {
		if token == numberToken {
			polishEntry = append(polishEntry, token)
		} else if token == openingBracket {
			stack.Push(token)
		} else if token == closingBracket {
			if stack.IsEmpty() {
				return nil, errBracketCombination
			}
			for {
				top, exists := stack.Peek()
				if !exists {
					return nil, errBracketCombination
				}
				if top == openingBracket {
					stack.Pop()
					break
				}
				last, _ := stack.Pop()
				polishEntry = append(polishEntry, last)
			}
		} else {
			for !stack.IsEmpty() {
				top, exists := stack.Peek()
				if !exists {
					break
				}
				if priority[top] >= priority[token] {
					last, _ := stack.Pop()
					polishEntry = append(polishEntry, last)
				} else {
					break
				}
			}
			stack.Push(token)
		}
	}

	for !stack.IsEmpty() {
		last, _ := stack.Pop()
		if last == openingBracket {
			return nil, errBracketCombination
		}
		polishEntry = append(polishEntry, last)
	}

	if len(polishEntry) == 0 {
		return nil, errInvalidFormula
	}

	return polishEntry, nil
}

// calculatePostfix обрабатывает постфиксную формулу, на выходе посчитанный результат.
func calculatePostfix(polishEntry []rune, numbers []float64) (float64, error) {
	stack := stack.New[float64]() // Используется для хранения чисел в определенной последовательности перед выполнением арифметических операций.
	numIndex := 0
	for _, token := range polishEntry {
		if token == numberToken {
			if numIndex >= len(numbers) {
				return 0, errors.New("invalid expression: not enough numbers")
			}
			stack.Push(numbers[numIndex])
			numIndex++
		} else {
			right, _ := stack.Pop()
			left, notEmpty := stack.Pop()
			if !notEmpty {
				return 0, errors.New("invalid expression: not enough values in stack")
			}

			result, err := operation(left, right, token)
			if err != nil {
				return 0, err
			}
			stack.Push(result)
		}
	}

	if stack.Size() != 1 {
		return 0, errors.New("invalid expression: multiple values left in stack")
	}

	resulst, _ := stack.Pop()

	return resulst, nil
}

// Calc высчитывает значение формулы в формате строки, возвращает посчитанный результат.
func Calc(formula string) (float64, error) {
	if err := isValidFormula(formula); err != nil {
		return 0, fmt.Errorf("calculate postfix: %v", err)
	}
	formula = addUnaryMinus(formula)

	tokens, numbers, err := tokenize(formula)
	if err != nil {
		return 0, fmt.Errorf("calculate postfix: %v", err)
	}

	polishEntry, err := infixToPostfix(tokens)
	if err != nil {
		return 0, fmt.Errorf("calculate postfix: %v", err)
	}

	result, err := calculatePostfix(polishEntry, numbers)
	if err != nil {
		return 0, fmt.Errorf("calculate postfix: %v", err)
	}

	return result, nil
}
