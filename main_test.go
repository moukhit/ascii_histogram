package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var blackhole map[byte]int

func BenchmarkProcessFolder(b *testing.B) {
	var dir string = "./data"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal("указанная папка не существует")
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("ошибка при получении списка файлов в указанной папке")
	}

	for i := 0; i < b.N; i++ {
		result := processFolder(dir, files)
		blackhole = result
	}
}
