package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Generator[T any] func(yield func(T) bool)

func WithOpen(filePath string, fn func(file *os.File) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("failed to close file: %v", err)
		}
	}(file)

	return fn(file)
}

func ReadJsonLineAsStruct[T any](filePath string, model *T) Generator[*T] {

	return func(yield func(*T) bool) {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatalf("failed to close file: %v", err)
			}
		}(file)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			err := json.Unmarshal([]byte(line), &model)

			if err != nil {
				log.Fatalf("error parsing JSON: %v", err)
			}

			if !yield(model) {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func WriteStringAsLine(filePath string, stringGenerator Generator[[]byte]) error {
	file, err := os.Create(filePath)
	if err != nil {
		return errors.New("error creating file: " + err.Error())
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("failed to close file: %v", err)
		}
	}(file)

	for val := range stringGenerator {
		_, err := file.Write(val)
		if err != nil {
			return errors.New("error writing file: " + err.Error())
		}

		_, err = file.WriteString("\n")
		if err != nil {
			return errors.New("error writing file: " + err.Error())
		}
	}

	err = file.Sync()
	if err != nil {
		return errors.New("error syncing file: " + err.Error())
	}
	return nil
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}
