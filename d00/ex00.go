package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"sort"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var numbers []int

	for scanner.Scan() {
		input := scanner.Text()

		// если введен пустой текст, то конец ввода
		if input == "" {
			break
		}

		num, err := strconv.Atoi(input)
		if err != nil {
			color.Red("На ввод должны поступать только числа! Ошибка: %v", err)
			break
		}

		if num >= 100000 || num <= -100000 {
			color.Red("Число выходит из диапазона [-100000, 100000].")
			break
		}

		numbers = append(numbers, num)
	}

	if err := scanner.Err(); err != nil {
		color.Red("Ошибка при чтении ввода:", err)
	}

	fmt.Println("Считанные числа:", numbers)

	// среднее значение
	if len(numbers) > 0 {
		sum := 0
		for _, num := range numbers {
			sum += num
		}
		average := float64(sum) / float64(len(numbers))

		fmt.Printf("Mean: %.2f\n", average)
	} else {
		fmt.Printf("Mean: %.2f\n", numbers[0])
	}

	// Вычисляем медиану
	if len(numbers) > 0 {
		sort.Ints(numbers)
		median := 0.0
		length := len(numbers)
		if length%2 == 1 {
			median = float64(numbers[length/2])
		} else {
			median = float64(numbers[length/2-1]+numbers[length/2]) / 2.0
		}

		fmt.Printf("Median: %.2f\n", median)
	} else {
		fmt.Printf("Median: %.2f\n", numbers[0])
	}

	// Вычисляем моду
	if len(numbers) > 0 {
		sort.Ints(numbers)
		mode := calculateMode(numbers)
		fmt.Printf("Mode: %d\n", mode)
	} else {
		fmt.Printf("Mode: %d\n", numbers[0])
	}

	fmt.Printf("SD: %.2f\n", standardDeviation(numbers))
}

// вычисление моды
func calculateMode(numbers []int) int {
	frequency := make(map[int]int)

	// частота встречаемости чисел
	for _, num := range numbers {
		frequency[num]++
	}

	// число с наибольшей частотой встречаемости
	maxFrequency := 0
	mode := 0
	for num, freq := range frequency {
		if freq > maxFrequency {
			maxFrequency = freq
			mode = num
		} else if freq == maxFrequency && num < mode {
			// если есть несколько чисел с одинаковой частотой, выбираем наименьшее из них
			mode = num
		}
	}

	return mode
}

func standardDeviation(numbers []int) float64 {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	mean := float64(sum) / float64(len(numbers))

	variance := 0.0

	for _, num := range numbers {
		diff := float64(num) - mean
		variance += diff * diff
	}

	variance /= float64(len(numbers))
	return (variance)
}
