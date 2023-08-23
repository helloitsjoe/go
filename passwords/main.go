package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
)

func md5Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func main() {
	data, err := os.ReadFile("wordlist.txt")
	if err != nil {
		fmt.Println("File reading error", err)
	}

	words := make([]string, 0)
	for _, word := range strings.Split(string(data), "\n") {
		words = append(words, string(word))
	}

	done := make(chan bool)

	for i := 0; i < 8; i++ {
		go func(i int) {
			offset := len(words) / 8
			start := i * offset
			end := start + offset

			for _, word1 := range words[start:end] {
				fmt.Println(word1)
				for _, word2 := range words {
					if md5Hash(word1+word2) == "3df7c3057e540cbe9244561a2d4345f7" {
						fmt.Println(word1, word2)
						done <- true
					}
				}
			}
		}(i)
	}

	isDone := <-done
	fmt.Println("Done", isDone)
	return
}
