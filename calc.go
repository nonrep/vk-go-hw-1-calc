package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Выполнение арифметической операции
func operation(left, right float64, operation rune) (result float64, err error) {
	switch operation {
	case '+':
		result = left + right
	case '-':
		result = left - right
	case '*':
		result = left * right
	case '/':
		if right == 0 {
			return result, errors.New("division by zero")
		}
		result = left / right
	}

	return result, nil
}

// Для сдвига слайса после выполнения арифметической операций
// index - перед каким элементом начать сдвиг
// skip - сколько элементов удалить перед index
func shiftSlice[T any](slice []T, index int, skip int) ([]T, error) {
	if index < 0 || index >= len(slice) {
		return slice, errors.New("index out of range")
	}
	indexWithSkip := index - skip
	if indexWithSkip < 0 {
		return slice, errors.New("index after skip is out of range")
	}
	result := make([]T, 0, len(slice)-skip)
	result = append(result, slice[:indexWithSkip]...)
	result = append(result, slice[index:]...)

	return result, nil
}

func calc(formula string, output io.Writer) (err error) {
	matched, err := regexp.MatchString(`^[\d\s+\-*/().]+$`, formula)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid characters")
	}

	// добавление поддержки унарного минуса путем добавления перед ним нуля
	var prev rune
	for i := range formula {
		if formula[i] == '-' && (prev == 0 || (prev >= 0 && prev <= 9) || prev == '(') {
			formula = formula[:i] + "0" + formula[i:]
			i++
		}
		if formula[i] != ' ' {
			prev = rune(formula[i])
		}
	}

	numberScanner := regexp.MustCompile(`\d+(\.\d+)?`)
	tokens := numberScanner.ReplaceAllString(formula, "N") // замена чисел на токен N для удобства обработки строки токенов
	if strings.Contains(tokens, ".") {
		return errors.New("the formula contains an invalid combination")
	}
	n := strings.Count(tokens, "N")
	stringNumbers := numberScanner.FindAllString(formula, n)
	tokens = strings.ReplaceAll(tokens, " ", "")
	var numbers []float64 // хранит все числа выражения чтобы сопоставить их с токенами чисел
	var number float64
	for _, string := range stringNumbers {
		number, err = strconv.ParseFloat(string, 64)
		if err != nil {
			return err
		}
		numbers = append(numbers, number)
	}
	if strings.Contains(tokens, "NN") {
		return errors.New("the formula contains an invalid combination")
	}

	// определение приоритета операций с помощью обратной польской записи
	priority := map[rune]int{
		'(': 0,
		')': 1,
		'+': 7,
		'-': 7,
		'/': 8,
		'*': 8,
	}
	var polishEntry []rune
	var stack []rune
	for _, token := range tokens {
		if token == 'N' {
			polishEntry = append(polishEntry, token)
		} else if token == '(' {
			stack = append(stack, token)
		} else if token == ')' {
			if len(stack) == 0 {
				return errors.New("wrong combination of brackets")
			}
			for stack[len(stack)-1] != '(' {
				polishEntry = append(polishEntry, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
		} else {
			for len(stack) > 0 && priority[stack[len(stack)-1]] >= priority[token] {
				polishEntry = append(polishEntry, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}
	for len(stack) != 0 {
		if stack[len(stack)-1] == '(' {
			return errors.New("wrong combination of brackets")
		}
		polishEntry = append(polishEntry, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	if len(polishEntry) == 0 {
		return errors.New("the formula contains an invalid combination")
	}

	// преобразование обратной польской записи в выражение
	// i - итератор polishEntry
	// num - итератор numbers
	var i int
	num := -1
	result := numbers[0] // если формула состоит из одного числа
	for len(polishEntry) != 1 {
		if polishEntry[i] == 'N' {
			i++
			num++
			if i > len(polishEntry)-1 {
				return errors.New("the formula contains an invalid combination")
			}
			continue
		} else {
			if i-2 >= 0 && num-1 >= 0 {
				result, err = operation(numbers[num-1], numbers[num], polishEntry[i])
				if err != nil {
					return err
				}
				numbers[num] = result
				polishEntry[i] = 'N'
				numbers, err = shiftSlice(numbers, num, 1)
				if err != nil {
					return err
				}
				polishEntry, err = shiftSlice(polishEntry, i, 2)
				if err != nil {
					return err
				}
				num -= 1
				i -= 1
			} else {
				return errors.New("the formula contains an invalid combination")
			}
		}
	}

	fmt.Fprintln(output, result)

	return nil
}

func main() {
	formula := os.Args[1]
	err := calc(formula, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
