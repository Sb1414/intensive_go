package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()

		// если введен пустой текст, то конец ввода
		if input == "" {
			break
		}

		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("На ввод должны поступать только числа!", err)
			break
		}

		if num >= 100000 || num <= -100000 {
			fmt.Println("Число выходит из диапазона [-100000, 100000].")
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
	}
}
