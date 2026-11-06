package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// Получаем текст со стандартного ввода
	text := ""
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

	// Разбиваем на слова
	var words []string
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

	// Заполняем массивы
	var names []string
	var ages []int
	var marks []int

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

	// Выполняем поиск по массивам

	// # 1. Средний балл у студентов, возраст <19
	avg := 0.0
	count := 0.0
	for i := range names {
		if ages[i] < 19 {
			avg += float64(marks[i])
			count++
		}
	}
	avg /= count
	fmt.Printf("> %f\n", avg)

	// # 2. Средний балл у студентов, имя на 'A'
	// avg := 0.0
	// count := 0.0
	// for i := range names {
	// 	// Кирилическая 'A' занимает 2 байта в unicode
	// 	if names[i][:2] == "А" {
	// 		avg += float64(marks[i])
	// 		count++
	// 	}
	// }
	// avg /= count
	// fmt.Printf("> %f\n", avg)

	// # 3. Максимальный балл у студентов, возраст <19, балл >70
	// max := -1
	// for i := range names {
	// 	if ages[i] < 19 && marks[i] > 70 {
	// 		if marks[i] > max {
	// 			max = marks[i]
	// 		}
	// 	}
	// }
	// fmt.Printf("> %d\n", max)
}
