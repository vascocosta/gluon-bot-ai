package utils

import (
	"bufio"
	"log"
	"os"
)

type ErrorAction string

const (
	Exit  ErrorAction = "exit"
	Fatal ErrorAction = "fatal"
)

func CheckError(err error, errAction ErrorAction) {
	if err != nil {
		switch errAction {
		case Exit:
			log.Println(err)
			os.Exit(1)
		case Fatal:
			log.Fatal(err)
		default:
			log.Fatal(err)
		}
	}
}

func Map[T any, U any](input []T, f func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = f(v)
	}

	return result
}

func ReadLastNLines(path string, n int) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
