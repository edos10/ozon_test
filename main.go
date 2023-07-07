package main

import (
	"fmt"
	"github.com"
)

type URL struct {
}

const sizeUrl = 10

func MakeMaps() map[int]string {
	result := make(map[int]string)
	c := rune(97)
	fmt.Printf(string(c))
	return result
}

func NextUrlString(current string) string {
	newUrl := []byte(current)
	for i := sizeUrl - 1; i > -1; i-- {
		if newUrl[]
	}
	result := string(newUrl)
	return result
}

func main() {
	c := "aaaaaaaaaa"
	fmt.Println(NextUrlString(c))
}
