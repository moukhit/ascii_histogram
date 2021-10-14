package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) != 2 {
		var sb strings.Builder
		sb.WriteString("неверное число переданных параметров\n")
		sb.WriteString("формат команды: ascii_hist <путь к папке с файлами>")
		log.Fatal(sb.String())
	}

	var dir string = os.Args[1]
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal("указанная папка не существует")
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("ошибка при получении списка файлов в указанной папке")
	}

	result := processFolder(dir, files)
	if result == nil {
		log.Println("отсутствуют файлы в указанной папке")
		return
	}

	printResult(result)
}

func processFolder(dir string, files []fs.FileInfo) map[byte]int {
	numFiles := len(files)
	if numFiles == 0 {
		return nil
	}

	maps := make(chan map[byte]int, numFiles)
	var wg sync.WaitGroup

	for _, f := range files {
		if !f.IsDir() {
			filename := filepath.Join(dir, f.Name())
			data, err := os.ReadFile(filename)
			if err != nil {
				continue
			}

			wg.Add(1)
			go func(data []byte) {
				defer wg.Done()

				m := map[byte]int{}
				for _, b := range data {
					m[b]++
				}
				maps <- m
			}(data)
		}
	}

	go func() {
		wg.Wait()
		close(maps)
	}()

	result := map[byte]int{}
	for m := range maps {
		for b, v := range m {
			result[b] += v
		}
	}

	return result
}

func printResult(result map[byte]int) {
	for k, v := range result {
		fmt.Println(string(k), ":", v)
	}
}
