package calc

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/nonrep/go-homework-1-calc/stack"
)

const numberToken = 'N' // Функция tokenize преобразует формулу в набор токенов, где все числа заменяются токеном N.

const plus = '+'
const minus = '-'
const multiplication = '*'
const division = '/'

var priority = map[rune]int{
	'(': 0,
	')': 1,
	'+': 7,
	'-': 7,
	'/': 8,
	'*': 8,
}

// operation выполняет арифметической операции.
func operation(left, right float64, operation rune) (result float64, err error) {
	switch operation {
	case plus:
		result = left + right
	case minus:
		result = left - right
	case multiplication:
		result = left * right
	case division:
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
		if char == minus && (prev == 0 || prev == ' ' || prev == '(') {
			formula = formula[:i] + "0" + formula[i:] // Добавляем 0 перед унарным минусом.
		}
		if char != ' ' {
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

		if unicode.IsDigit(char) || char == '.' {
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

			if char == plus || char == minus || char == multiplication || char == division || char == '(' || char == ')' {
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
		return []rune{}, errors.New("the formula contains an invalid combination")
	}

	var polishEntry []rune

	stack := stack.Stack[rune]{}

	for _, token := range tokens {
		if token == numberToken {
			polishEntry = append(polishEntry, token)
		} else if token == '(' {
			stack.Push(token)
		} else if token == ')' {
			if stack.IsEmpty() {
				return nil, errors.New("wrong combination of brackets")
			}
			for {
				top, exists := stack.Peek()
				if !exists {
					return nil, errors.New("wrong combination of brackets")
				}
				if top == '(' {
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
		if last == '(' {
			return nil, errors.New("wrong combination of brackets")
		}
		polishEntry = append(polishEntry, last)
	}

	if len(polishEntry) == 0 {
		return nil, errors.New("the formula contains an invalid combination")
	}

	return polishEntry, nil
}

// calculatePostfix обрабатывает постфиксную формулу, на выходе посчитанный результат.
func calculatePostfix(polishEntry []rune, numbers []float64) (float64, error) {
	var stack []float64
	numIndex := 0

	for _, token := range polishEntry {
		if token == numberToken {
			if numIndex >= len(numbers) {
				return 0, errors.New("invalid expression: not enough numbers")
			}
			stack = append(stack, numbers[numIndex])
			numIndex++
		} else {
			if len(stack) < 2 {
				return 0, errors.New("invalid expression: not enough values in stack")
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			result, err := operation(left, right, token)
			if err != nil {
				return 0, err
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression: multiple values left in stack")
	}

	return stack[0], nil
}

// Calc высчитывает значение формулы в формате строки, возвращает посчитанный результат.
func Calc(formula string) (float64, error) {
	if err := isValidFormula(formula); err != nil {
		return 0, err
	}
	formula = addUnaryMinus(formula)

	tokens, numbers, err := tokenize(formula)
	if err != nil {
		return 0, err
	}

	polishEntry, err := infixToPostfix(tokens)
	if err != nil {
		return 0, err
	}

	return calculatePostfix(polishEntry, numbers)
}
