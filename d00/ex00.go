package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
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
}
