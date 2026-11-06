package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	text  string
	words []string

	names []string
	ages  []int
	marks []int

	findResult []int

	result float64
)

func main() {
	ReadAllText()
	SplitTextByWords()
	SplitWordsToArrays()

	// # 1. Средний балл у студентов, возраст <19
	FindStudentsSubAge(19)
	ResultAvgMarkForList(findResult)

	// # 2. Средний балл у студентов, имя на 'A'
	// FindStudentsNameStartsWith("А")
	// ResultAvgMarkForList(findResult)

	// # 3. Максимальный балл у студентов, возраст <19, балл >70
	// FindStudentsSubAgeAndOverMark(19, 70)
	// ResultMaxMarkForList(findResult)

	PrintResult()
}

func ReadAllText() {
	buff := make([]byte, 1024)
	for {
		n, err := os.Stdin.Read(buff)
		if err == io.EOF || n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			panic(err)
		}

		text += string(buff[:n])
	}
}

func SplitTextByWords() {
	word := ""
	for _, c := range text {
		if c == ' ' || c == '\n' {
			words = append(words, word)
			word = ""
		} else {
			word += string(c)
		}
	}
	if word != "" {
		words = append(words, word)
	}
}

func SplitWordsToArrays() {
	if len(words)%3 != 0 {
		panic("invalid arguments count")
	}
	for i := 0; i < len(words); i += 3 {
		// Имена
		name := words[i]
		names = append(names, name)

		// Возрасты
		age := words[i+1]
		num := 0
		exp := 1
		// придётся использовать несколько встроенных функций
		for i := len(age) - 1; i >= 0; i-- {
			// Преобразуем числа по цифрам из byte (char)
			if age[i] < '0' || age[i] > '9' {
				panic(age + " bad value for Age")
			}
			n := int(age[i] - '0')

			num += n * exp
			exp *= 10
		}
		ages = append(ages, num)

		// Оценки
		mark := words[i+2]
		num = 0
		exp = 1
		for i := len(mark) - 1; i >= 0; i-- {
			// Преобразуем числа по цифрам из byte (char)
			if mark[i] < '0' || mark[i] > '9' {
				panic(mark + " bad value for Age")
			}
			n := int(mark[i] - '0')

			num += n * exp
			exp *= 10
		}
		marks = append(marks, num)
	}
}

func FindStudentsSubAge(age int) {
	findResult = []int{}
	for i := range names {
		if ages[i] < age {
			findResult = append(findResult, i)
		}
	}
}

func FindStudentsSubAgeAndOverMark(age int, mark int) {
	findResult = []int{}
	for i := range names {
		if ages[i] < age && marks[i] > mark {
			findResult = append(findResult, i)
		}
	}
}

func FindStudentsNameStartsWith(c string) {
	findResult = []int{}
	for i := range names {
		if strings.HasPrefix(names[i], c) {
			findResult = append(findResult, i)
		}
	}
}

func ResultAvgMarkForList(idxs []int) {
	avg := 0.0
	count := 0.0
	for _, i := range idxs {
		avg += float64(marks[i])
		count++
	}

	result = avg / count
}

func ResultMaxMarkForList(idxs []int) {
	max := -1.0
	for _, i := range idxs {
		if float64(marks[i]) > max {
			max = float64(marks[i])
		}
	}

	result = max
}

func PrintResult() {
	fmt.Printf("> %f\n", result)
}
