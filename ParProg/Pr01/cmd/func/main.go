package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	TEXT := GetInputText(os.Stdin)
	WORDS := SplitByWords(TEXT)
	DATA := CombineWordsToData(WORDS)

	// # 1. Средний балл у студентов, возраст <19
	FILTRED_DATA := Find(DATA, func(d Data, idx int) bool {
		return d.Ages[idx] < 19
	})
	PrintFormat(GetAverage(FILTRED_DATA.Marks))

	// # 2. Средний балл у студентов, имя на 'A'
	// FILTRED_DATA := Find(DATA, func(d Data, idx int) bool {
	// 	return strings.HasPrefix(d.Names[idx], "А")
	// })
	// PrintFormat(GetAverage(FILTRED_DATA.Marks))

	// # 3. Максимальный балл у студентов, возраст <19, балл >70
	// FILTRED_DATA := Find(DATA, func(d Data, idx int) bool {
	// 	return d.Ages[idx] < 19 && d.Marks[idx] > 70
	// })
	// PrintFormat(GetMax(FILTRED_DATA.Marks))

}

type Data struct {
	Names []string
	Ages  []int
	Marks []int
}

func GetInputText(src *os.File) (text string) {
	buff := make([]byte, 1024)
	for {
		n, err := src.Read(buff)
		if err == io.EOF || n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			panic(err)
		}

		text += string(buff[:n])
	}
	return text
}

func SplitByWords(text string) (words []string) {
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
	return words
}

func atoi(s string) int {
	num := 0
	exp := 1
	for i := len(s) - 1; i >= 0; i-- {
		// Преобразуем числа по цифрам из byte (char)
		if s[i] < '0' || s[i] > '9' {
			panic(s + " bad value for Age")
		}
		n := int(s[i] - '0')

		num += n * exp
		exp *= 10
	}
	return num
}

func CombineWordsToData(words []string) Data {
	if len(words)%3 != 0 {
		panic("invalid arguments count")
	}
	var (
		names []string
		ages  []int
		marks []int
	)
	for i := 0; i < len(words); i += 3 {
		// Имена
		name := words[i]
		names = append(names, name)

		// Возрасты
		age := words[i+1]
		ages = append(ages, atoi(age))

		// Оценки
		mark := words[i+2]
		marks = append(marks, atoi(mark))
	}
	return Data{
		Names: names,
		Ages:  ages,
		Marks: marks,
	}
}

type predicate func(d Data, idx int) bool

func Find(d Data, f predicate) Data {
	var (
		names []string
		ages  []int
		marks []int
	)
	for i := range d.Names {
		if f(d, i) {
			names = append(names, d.Names[i])
			ages = append(ages, d.Ages[i])
			marks = append(marks, d.Marks[i])
		}
	}
	return Data{
		Names: names,
		Ages:  ages,
		Marks: marks,
	}
}

func GetAverage(nums []int) float64 {
	avg := 0.0
	count := 0.0
	for _, n := range nums {
		avg += float64(n)
		count++
	}

	return avg / count
}

func GetMax(nums []int) float64 {
	max := -1.0
	for _, n := range nums {
		if float64(n) > max {
			max = float64(n)
		}
	}

	return max
}

func PrintFormat(n any) {
	fmt.Println("> ", n)
}
